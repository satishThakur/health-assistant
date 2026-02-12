# TimescaleDB Aggregation Strategy

## Overview

This document outlines when and how to implement TimescaleDB's continuous aggregates and time-bucketing features for the Health Assistant application. It serves as a future reference for optimizing query performance as data volume grows.

## Current State (2026-02-11)

### Data Granularity
- **Sleep data**: 1 event per night (daily)
- **Daily stats**: 1 event per day (steps, calories, distance, heart rate)
- **Activities**: 1 event per session (workouts)
- **HRV**: 1 event per day
- **Stress**: 1 event per day (or daily summary)
- **Body Battery**: 1 event per day

### Query Patterns
- Dashboard queries: Last 24 hours (fast)
- Week trends: Last 7 days (fast)
- Correlations: Last 30-90 days (acceptable)

### Performance Baseline
- Event count: ~10-20 events/day per user
- Query latency: <100ms for all current queries
- Database size: Minimal (days/weeks of data)

## Implementation Triggers

### When to Implement Continuous Aggregates

Implement when **ANY** of these conditions are met:

1. **Query Latency Threshold**
   - Week trends query >500ms
   - Correlation queries >1 second
   - Dashboard load >200ms

2. **Data Volume Threshold**
   - >90 days of historical data per user
   - >100,000 events in database
   - Multiple users with long histories

3. **New High-Frequency Data Streams**
   - Adding minute-by-minute heart rate
   - Continuous stress readings (5-min intervals)
   - Real-time step counting
   - Intraday body battery updates

4. **Complex Aggregation Requirements**
   - Moving averages (7-day, 30-day, 90-day)
   - Percentile calculations across large datasets
   - Multi-metric correlations requiring joins

### Pre-Implementation Checklist

Before implementing, confirm:
- [ ] Queries are actually slow (measure with EXPLAIN ANALYZE)
- [ ] Raw data is needed for at least some queries
- [ ] You understand what granularity users need
- [ ] You have a data retention strategy
- [ ] You've tested on production-like data volume

## Architecture Design

### Three-Tier Aggregation Hierarchy

```
Raw Events (hypertable)
    ↓
Daily Aggregates (continuous aggregate)
    ↓
Weekly/Monthly Aggregates (continuous aggregate from daily)
```

### Data Flow

```
Garmin Scheduler → Ingestion Service → events (hypertable)
                                            ↓
                          daily_health_stats (continuous aggregate)
                                            ↓
                          weekly_health_stats (continuous aggregate)
                                            ↓
                          monthly_health_stats (continuous aggregate)
```

## Implementation Guide

### Step 1: Analyze Current Query Performance

```sql
-- Enable query timing
\timing on

-- Test current week trends query
EXPLAIN ANALYZE
SELECT time, event_type, data
FROM events
WHERE user_id = 'test-user-id'
  AND time >= NOW() - INTERVAL '7 days'
ORDER BY time ASC;

-- Test correlation query
EXPLAIN ANALYZE
SELECT time, event_type, data
FROM events
WHERE user_id = 'test-user-id'
  AND time >= NOW() - INTERVAL '30 days'
ORDER BY time ASC;
```

**Decision Point**: If queries take <500ms, **don't implement yet**.

### Step 2: Create Daily Health Stats Aggregate

For when you have minute-level data or need fast daily rollups:

```sql
-- Daily aggregate of all health metrics
CREATE MATERIALIZED VIEW daily_health_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) AS day,
    user_id,
    event_type,

    -- Sleep metrics (from garmin_sleep events)
    AVG(NULLIF((data->>'duration_minutes')::numeric, 0)) AS avg_sleep_minutes,
    AVG(NULLIF((data->>'sleep_score')::numeric, 0)) AS avg_sleep_score,
    AVG(NULLIF((data->>'hrv_avg')::numeric, 0)) AS avg_hrv_during_sleep,

    -- Activity metrics (from garmin_daily_stats events)
    SUM((data->>'steps')::numeric) AS total_steps,
    SUM((data->>'calories')::numeric) AS total_calories,
    SUM((data->>'distance_meters')::numeric) AS total_distance_meters,
    AVG(NULLIF((data->>'resting_heart_rate')::numeric, 0)) AS avg_resting_hr,
    MAX((data->>'max_heart_rate')::numeric) AS max_heart_rate,

    -- Stress metrics (from garmin_stress events)
    AVG(NULLIF((data->>'average')::numeric, 0)) AS avg_stress_level,

    -- Body battery metrics (from garmin_body_battery events)
    SUM((data->>'charged')::numeric) AS total_charged,
    SUM((data->>'drained')::numeric) AS total_drained,

    -- Activity sessions count
    COUNT(*) FILTER (WHERE event_type = 'garmin_activity') AS activity_session_count,
    SUM((data->>'duration_minutes')::numeric) FILTER (WHERE event_type = 'garmin_activity') AS total_activity_minutes,

    COUNT(*) AS event_count
FROM events
GROUP BY day, user_id, event_type
WITH NO DATA;

-- Refresh policy: Update last 3 days every hour
SELECT add_continuous_aggregate_policy('daily_health_stats',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- Initial data population
CALL refresh_continuous_aggregate('daily_health_stats', NULL, NULL);
```

### Step 3: Create Weekly Aggregates

```sql
-- Weekly rollup from daily stats
CREATE MATERIALIZED VIEW weekly_health_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('7 days', day) AS week,
    user_id,

    -- Sleep averages
    AVG(avg_sleep_minutes) AS avg_sleep_minutes,
    AVG(avg_sleep_score) AS avg_sleep_score,
    AVG(avg_hrv_during_sleep) AS avg_hrv,

    -- Activity totals and averages
    SUM(total_steps) AS total_steps,
    AVG(total_steps) AS avg_daily_steps,
    SUM(total_calories) AS total_calories,
    SUM(total_distance_meters) AS total_distance_meters,
    AVG(avg_resting_hr) AS avg_resting_hr,
    MAX(max_heart_rate) AS max_heart_rate,

    -- Stress
    AVG(avg_stress_level) AS avg_stress_level,

    -- Body battery
    SUM(total_charged) AS total_charged,
    SUM(total_drained) AS total_drained,

    -- Activity sessions
    SUM(activity_session_count) AS total_activity_sessions,
    SUM(total_activity_minutes) AS total_activity_minutes,
    AVG(total_activity_minutes) AS avg_daily_activity_minutes,

    COUNT(DISTINCT day) AS days_with_data
FROM daily_health_stats
GROUP BY week, user_id
WITH NO DATA;

-- Refresh policy: Update weekly stats daily
SELECT add_continuous_aggregate_policy('weekly_health_stats',
    start_offset => INTERVAL '14 days',
    end_offset => INTERVAL '1 day',
    schedule_interval => INTERVAL '1 day');

-- Initial population
CALL refresh_continuous_aggregate('weekly_health_stats', NULL, NULL);
```

### Step 4: Create Monthly Aggregates

```sql
-- Monthly rollup for long-term trends
CREATE MATERIALIZED VIEW monthly_health_stats
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('30 days', week) AS month,
    user_id,

    -- All metrics as averages or totals
    AVG(avg_sleep_minutes) AS avg_sleep_minutes,
    AVG(avg_sleep_score) AS avg_sleep_score,
    AVG(avg_hrv) AS avg_hrv,

    SUM(total_steps) AS total_steps,
    AVG(avg_daily_steps) AS avg_daily_steps,
    SUM(total_calories) AS total_calories,
    AVG(avg_resting_hr) AS avg_resting_hr,

    AVG(avg_stress_level) AS avg_stress_level,

    SUM(total_activity_sessions) AS total_activity_sessions,
    SUM(total_activity_minutes) AS total_activity_minutes,

    SUM(days_with_data) AS days_with_data
FROM weekly_health_stats
GROUP BY month, user_id
WITH NO DATA;

-- Refresh policy: Update monthly stats weekly
SELECT add_continuous_aggregate_policy('monthly_health_stats',
    start_offset => INTERVAL '60 days',
    end_offset => INTERVAL '7 days',
    schedule_interval => INTERVAL '7 days');

-- Initial population
CALL refresh_continuous_aggregate('monthly_health_stats', NULL, NULL);
```

### Step 5: Update Go Backend Queries

Create new repository methods to leverage aggregates:

```go
// backend/internal/db/aggregates.go

type AggregateRepository struct {
    db *Database
}

func NewAggregateRepository(db *Database) *AggregateRepository {
    return &AggregateRepository{db: db}
}

// GetWeeklyStats retrieves pre-aggregated weekly statistics
func (r *AggregateRepository) GetWeeklyStats(ctx context.Context, userID string, weeks int) ([]WeeklyStats, error) {
    query := `
        SELECT
            week,
            avg_sleep_minutes,
            avg_sleep_score,
            total_steps,
            avg_daily_steps,
            total_calories,
            avg_resting_hr,
            avg_stress_level,
            total_activity_sessions,
            avg_daily_activity_minutes,
            days_with_data
        FROM weekly_health_stats
        WHERE user_id = $1
            AND week >= NOW() - INTERVAL '1 week' * $2
        ORDER BY week DESC
    `

    rows, err := r.db.Pool.Query(ctx, query, userID, weeks)
    if err != nil {
        return nil, fmt.Errorf("failed to query weekly stats: %w", err)
    }
    defer rows.Close()

    var stats []WeeklyStats
    for rows.Next() {
        var ws WeeklyStats
        if err := rows.Scan(
            &ws.Week,
            &ws.AvgSleepMinutes,
            &ws.AvgSleepScore,
            &ws.TotalSteps,
            &ws.AvgDailySteps,
            &ws.TotalCalories,
            &ws.AvgRestingHR,
            &ws.AvgStressLevel,
            &ws.TotalActivitySessions,
            &ws.AvgDailyActivityMinutes,
            &ws.DaysWithData,
        ); err != nil {
            return nil, err
        }
        stats = append(stats, ws)
    }

    return stats, nil
}

// GetMonthlyTrends retrieves long-term trend data
func (r *AggregateRepository) GetMonthlyTrends(ctx context.Context, userID string, months int) ([]MonthlyStats, error) {
    query := `
        SELECT
            month,
            avg_sleep_minutes,
            avg_sleep_score,
            total_steps,
            avg_daily_steps,
            total_calories,
            avg_resting_hr,
            avg_stress_level,
            total_activity_sessions,
            total_activity_minutes,
            days_with_data
        FROM monthly_health_stats
        WHERE user_id = $1
            AND month >= NOW() - INTERVAL '30 days' * $2
        ORDER BY month DESC
    `

    // Similar implementation
}
```

### Step 6: Add Compression and Retention Policies

```sql
-- Compress data older than 30 days
SELECT add_compression_policy('events', INTERVAL '30 days');

-- Keep raw events for 1 year (aggregates remain)
-- WARNING: Only enable after confirming aggregates work correctly
-- SELECT add_retention_policy('events', INTERVAL '365 days');

-- Keep daily aggregates for 2 years
SELECT add_retention_policy('daily_health_stats', INTERVAL '730 days');

-- Keep weekly aggregates for 5 years
SELECT add_retention_policy('weekly_health_stats', INTERVAL '1825 days');

-- Keep monthly aggregates forever (or very long)
-- No retention policy on monthly_health_stats
```

## Migration Strategy

### Phase 1: Create Aggregates (Non-Disruptive)
1. Create continuous aggregate views
2. Let them populate in background
3. Monitor refresh times and storage impact
4. **Keep existing queries unchanged**

### Phase 2: A/B Testing
1. Add new endpoints using aggregates (e.g., `/api/v1/trends/weekly-aggregated`)
2. Test performance vs old endpoints
3. Validate data accuracy
4. Compare query latency

### Phase 3: Gradual Migration
1. Update internal queries to use aggregates
2. Add feature flag to switch between raw and aggregated queries
3. Monitor for data discrepancies
4. Roll back if needed

### Phase 4: Cleanup
1. Remove old query methods
2. Enable compression on raw events
3. Consider retention policies after 6 months
4. Update documentation

## Performance Monitoring

### Queries to Track Performance

```sql
-- Check continuous aggregate size
SELECT
    hypertable_name,
    pg_size_pretty(total_bytes) AS total_size,
    pg_size_pretty(table_bytes) AS table_size
FROM timescaledb_information.hypertables
ORDER BY total_bytes DESC;

-- Check refresh job status
SELECT
    job_id,
    application_name,
    last_run_status,
    last_run_started_at,
    last_successful_finish,
    next_start,
    total_runs,
    total_successes,
    total_failures
FROM timescaledb_information.jobs
WHERE application_name LIKE '%continuous_aggregate%';

-- Check query performance on aggregates
EXPLAIN ANALYZE
SELECT * FROM weekly_health_stats
WHERE user_id = 'test-user-id'
ORDER BY week DESC
LIMIT 12;
```

### Alert Thresholds

Set up monitoring for:
- Aggregate refresh job failures
- Refresh duration >5 minutes
- Query latency regression (aggregates slower than raw)
- Storage growth rate

## Cost-Benefit Analysis

### Benefits
- **Query speed**: 10-100x faster for historical queries
- **Storage efficiency**: Compression reduces storage by 60-90%
- **Scalability**: Handles years of data without degradation
- **Analytics**: Enable complex correlations without expensive scans

### Costs
- **Complexity**: Additional views to maintain
- **Storage**: Aggregates use additional space (until raw data compressed/deleted)
- **Refresh overhead**: Background jobs consume CPU/IO
- **Development time**: Migration and testing effort

### When Benefits Outweigh Costs
- >100,000 events in database
- Queries taking >1 second
- Users requesting long-term trends (6+ months)
- Multiple users querying simultaneously
- Adding real-time dashboards

## Alternative: Time-Bucketing Without Continuous Aggregates

For moderate data volumes, use TimescaleDB's time_bucket in regular queries:

```sql
-- Weekly aggregation on-the-fly (no continuous aggregate needed)
SELECT
    time_bucket('7 days', time) AS week,
    AVG((data->>'steps')::numeric) AS avg_steps,
    SUM((data->>'calories')::numeric) AS total_calories
FROM events
WHERE user_id = $1
    AND event_type = 'garmin_daily_stats'
    AND time >= NOW() - INTERVAL '90 days'
GROUP BY week
ORDER BY week DESC;
```

**Use this approach if:**
- You have <90 days of data
- Queries are <500ms
- User base is small (<100 active users)

## Future Considerations

### Intraday Data Support

If adding minute-level data streams (e.g., continuous heart rate):

```sql
-- Minute-level heart rate readings
CREATE TABLE heart_rate_readings (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL,
    heart_rate INT NOT NULL,
    confidence DECIMAL(3,2)
);

-- Convert to hypertable
SELECT create_hypertable('heart_rate_readings', 'time');

-- 5-minute aggregate
CREATE MATERIALIZED VIEW heart_rate_5min
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('5 minutes', time) AS bucket,
    user_id,
    AVG(heart_rate) AS avg_hr,
    MIN(heart_rate) AS min_hr,
    MAX(heart_rate) AS max_hr,
    COUNT(*) AS reading_count
FROM heart_rate_readings
GROUP BY bucket, user_id;

-- Hourly aggregate from 5-min buckets
CREATE MATERIALIZED VIEW heart_rate_hourly
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', bucket) AS hour,
    user_id,
    AVG(avg_hr) AS avg_hr,
    MIN(min_hr) AS min_hr,
    MAX(max_hr) AS max_hr,
    SUM(reading_count) AS reading_count
FROM heart_rate_5min
GROUP BY hour, user_id;
```

### Real-Time Aggregates

For live dashboards requiring <1 minute latency:

```sql
-- Real-time refresh (higher CPU cost)
SELECT add_continuous_aggregate_policy('heart_rate_5min',
    start_offset => INTERVAL '30 minutes',
    end_offset => INTERVAL '1 minute',  -- Very recent data
    schedule_interval => INTERVAL '1 minute');  -- Refresh every minute
```

## References

- [TimescaleDB Continuous Aggregates Docs](https://docs.timescale.com/use-timescale/latest/continuous-aggregates/)
- [TimescaleDB Compression](https://docs.timescale.com/use-timescale/latest/compression/)
- [TimescaleDB Data Retention](https://docs.timescale.com/use-timescale/latest/data-retention/)
- [TimescaleDB Best Practices](https://docs.timescale.com/use-timescale/latest/best-practices/)

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-02-11 | Defer implementation | Data volume too small, queries fast enough |
| TBD | Implement daily aggregates | Query latency exceeded 500ms threshold |
| TBD | Enable compression | Database size exceeded 10GB |
| TBD | Add retention policy | After validating aggregates work for 3 months |

## Questions Before Implementation

1. What is the current slowest query and its latency?
2. How many total events are in the database?
3. What is the oldest event in the database?
4. Are there any queries that scan the entire events table?
5. What percentage of queries would benefit from aggregates?
6. Do we have monitoring in place to detect query regressions?
7. Have we tested aggregate refresh time with production data volume?
8. What is our disaster recovery plan if aggregates get out of sync?

## Conclusion

**Current recommendation: Wait and monitor.**

Implement continuous aggregates when:
1. You have >90 days of data, OR
2. Queries take >500ms, OR
3. You add high-frequency data streams (minute-level readings)

Until then, focus on:
- Completing core features
- Adding proper query monitoring
- Testing with realistic data volumes
- Documenting slow queries as they emerge

This document should be reviewed quarterly as the application scales.

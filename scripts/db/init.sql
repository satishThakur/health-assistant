-- Initial database schema for Health Assistant
-- PostgreSQL + TimescaleDB

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    garmin_oauth_token JSONB,
    preferences JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Events table (TimescaleDB hypertable)
CREATE TABLE IF NOT EXISTS events (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    source VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    metadata JSONB,
    confidence FLOAT CHECK (confidence >= 0 AND confidence <= 1),
    PRIMARY KEY (time, user_id, event_type)
);

-- Convert to hypertable
SELECT create_hypertable('events', 'time', if_not_exists => TRUE);

-- Indexes for events
CREATE INDEX IF NOT EXISTS idx_events_user_time ON events (user_id, time DESC);
CREATE INDEX IF NOT EXISTS idx_events_type_time ON events (event_type, time DESC);
CREATE INDEX IF NOT EXISTS idx_events_data_gin ON events USING GIN (data);

-- Experiments table
CREATE TABLE IF NOT EXISTS experiments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    hypothesis TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('proposed', 'accepted', 'active', 'completed', 'abandoned')),
    intervention JSONB NOT NULL,
    control_condition JSONB,
    duration_days INT NOT NULL,
    start_date DATE,
    end_date DATE,
    compliance_rate FLOAT,
    results JSONB,
    posterior_beliefs JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_experiments_user_status ON experiments (user_id, status);
CREATE INDEX IF NOT EXISTS idx_experiments_dates ON experiments (start_date, end_date);

-- Continuous aggregates for common queries (TimescaleDB feature)
-- Daily metrics rollup
CREATE MATERIALIZED VIEW IF NOT EXISTS daily_metrics
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) AS day,
    user_id,
    event_type,
    COUNT(*) as event_count,
    data
FROM events
GROUP BY day, user_id, event_type, data
WITH NO DATA;

-- Refresh policy for continuous aggregate
SELECT add_continuous_aggregate_policy('daily_metrics',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour',
    if_not_exists => TRUE);

-- Data retention policy (optional - keep data for 2 years)
-- SELECT add_retention_policy('events', INTERVAL '2 years', if_not_exists => TRUE);

-- Insert a test user for development
INSERT INTO users (id, email, password_hash, preferences)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'test@healthassistant.dev',
    '$2a$10$dummyhashfordevonly',  -- Not a real hash
    '{"timezone": "America/New_York", "subjective_reminders": ["08:00", "22:00"]}'::jsonb
)
ON CONFLICT (email) DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO healthuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO healthuser;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Health Assistant database initialized successfully!';
    RAISE NOTICE 'TimescaleDB extension enabled';
    RAISE NOTICE 'Test user created: test@healthassistant.dev';
END $$;

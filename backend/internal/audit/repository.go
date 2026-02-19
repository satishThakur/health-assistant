package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/satishthakur/health-assistant/backend/internal/db"
)

// SyncAudit represents a sync audit record.
type SyncAudit struct {
	ID                  string     `json:"id"`
	SyncStartedAt       time.Time  `json:"sync_started_at"`
	SyncCompletedAt     *time.Time `json:"sync_completed_at,omitempty"`
	SyncDurationSeconds *int       `json:"sync_duration_seconds,omitempty"`
	UserID              string     `json:"user_id"`
	DataType            string     `json:"data_type"`
	TargetDate          string     `json:"target_date"`
	RecordsFetched      int        `json:"records_fetched"`
	RecordsInserted     int        `json:"records_inserted"`
	RecordsUpdated      int        `json:"records_updated"`
	EarliestTimestamp   *time.Time `json:"earliest_timestamp,omitempty"`
	LatestTimestamp     *time.Time `json:"latest_timestamp,omitempty"`
	Status              string     `json:"status"`
	ErrorMessage        *string    `json:"error_message,omitempty"`
	Metadata            []byte     `json:"metadata,omitempty"`
}

// Repository handles database operations for sync audit.
type Repository struct {
	db *db.Database
}

// NewRepository creates a new audit Repository.
func NewRepository(database *db.Database) *Repository {
	return &Repository{db: database}
}

// InsertSyncAudit inserts a new sync audit record.
func (r *Repository) InsertSyncAudit(ctx context.Context, audit *SyncAudit) error {
	query := `
		INSERT INTO sync_audit (
			sync_started_at, sync_completed_at, sync_duration_seconds,
			user_id, data_type, target_date,
			records_fetched, records_inserted, records_updated,
			earliest_timestamp, latest_timestamp,
			status, error_message, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	var id string
	err := r.db.Pool.QueryRow(
		ctx, query,
		audit.SyncStartedAt, audit.SyncCompletedAt, audit.SyncDurationSeconds,
		audit.UserID, audit.DataType, audit.TargetDate,
		audit.RecordsFetched, audit.RecordsInserted, audit.RecordsUpdated,
		audit.EarliestTimestamp, audit.LatestTimestamp,
		audit.Status, audit.ErrorMessage, audit.Metadata,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert sync audit: %w", err)
	}

	audit.ID = id
	return nil
}

// GetRecentSyncAudits retrieves recent sync audit records.
func (r *Repository) GetRecentSyncAudits(ctx context.Context, userID string, limit int) ([]SyncAudit, error) {
	query := `
		SELECT
			id, sync_started_at, sync_completed_at, sync_duration_seconds,
			user_id, data_type, target_date,
			records_fetched, records_inserted, records_updated,
			earliest_timestamp, latest_timestamp,
			status, error_message, metadata
		FROM sync_audit
		WHERE user_id = $1
		ORDER BY sync_started_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query sync audits: %w", err)
	}
	defer rows.Close()

	var audits []SyncAudit
	for rows.Next() {
		var a SyncAudit
		if err := rows.Scan(
			&a.ID, &a.SyncStartedAt, &a.SyncCompletedAt, &a.SyncDurationSeconds,
			&a.UserID, &a.DataType, &a.TargetDate,
			&a.RecordsFetched, &a.RecordsInserted, &a.RecordsUpdated,
			&a.EarliestTimestamp, &a.LatestTimestamp,
			&a.Status, &a.ErrorMessage, &a.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sync audit: %w", err)
		}
		audits = append(audits, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sync audits: %w", err)
	}

	return audits, nil
}

// GetSyncAuditsByDataType retrieves sync audit records filtered by data type.
func (r *Repository) GetSyncAuditsByDataType(ctx context.Context, dataType string, limit int) ([]SyncAudit, error) {
	query := `
		SELECT
			id, sync_started_at, sync_completed_at, sync_duration_seconds,
			user_id, data_type, target_date,
			records_fetched, records_inserted, records_updated,
			earliest_timestamp, latest_timestamp,
			status, error_message, metadata
		FROM sync_audit
		WHERE data_type = $1
		ORDER BY sync_started_at DESC
		LIMIT $2
	`

	rows, err := r.db.Pool.Query(ctx, query, dataType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query sync audits: %w", err)
	}
	defer rows.Close()

	var audits []SyncAudit
	for rows.Next() {
		var a SyncAudit
		if err := rows.Scan(
			&a.ID, &a.SyncStartedAt, &a.SyncCompletedAt, &a.SyncDurationSeconds,
			&a.UserID, &a.DataType, &a.TargetDate,
			&a.RecordsFetched, &a.RecordsInserted, &a.RecordsUpdated,
			&a.EarliestTimestamp, &a.LatestTimestamp,
			&a.Status, &a.ErrorMessage, &a.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sync audit: %w", err)
		}
		audits = append(audits, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sync audits: %w", err)
	}

	return audits, nil
}

// GetSyncAuditStats retrieves summary statistics for sync audits.
func (r *Repository) GetSyncAuditStats(ctx context.Context, userID string, startDate, endDate time.Time) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total_syncs,
			SUM(records_fetched) as total_fetched,
			SUM(records_inserted) as total_inserted,
			SUM(records_updated) as total_updated,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_syncs,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_syncs,
			AVG(sync_duration_seconds) as avg_duration_seconds
		FROM sync_audit
		WHERE user_id = $1
			AND sync_started_at >= $2
			AND sync_started_at <= $3
	`

	var stats struct {
		TotalSyncs         int
		TotalFetched       int
		TotalInserted      int
		TotalUpdated       int
		SuccessfulSyncs    int
		FailedSyncs        int
		AvgDurationSeconds *float64
	}

	err := r.db.Pool.QueryRow(ctx, query, userID, startDate, endDate).Scan(
		&stats.TotalSyncs, &stats.TotalFetched, &stats.TotalInserted,
		&stats.TotalUpdated, &stats.SuccessfulSyncs, &stats.FailedSyncs,
		&stats.AvgDurationSeconds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync audit stats: %w", err)
	}

	result := map[string]interface{}{
		"total_syncs":      stats.TotalSyncs,
		"total_fetched":    stats.TotalFetched,
		"total_inserted":   stats.TotalInserted,
		"total_updated":    stats.TotalUpdated,
		"successful_syncs": stats.SuccessfulSyncs,
		"failed_syncs":     stats.FailedSyncs,
	}

	if stats.AvgDurationSeconds != nil {
		result["avg_duration_seconds"] = *stats.AvgDurationSeconds
	}

	return result, nil
}

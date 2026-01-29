package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/satishthakur/health-assistant/backend/internal/models"
)

// EventRepository handles database operations for events
type EventRepository struct {
	db *Database
}

// NewEventRepository creates a new EventRepository
func NewEventRepository(db *Database) *EventRepository {
	return &EventRepository{db: db}
}

// InsertEventResult contains the result of an insert/update operation
type InsertEventResult struct {
	WasInserted bool // true if inserted, false if updated
}

// InsertEvent inserts a new event or updates if conflict on (time, user_id, event_type)
// Returns InsertEventResult indicating whether the row was inserted or updated
func (r *EventRepository) InsertEvent(ctx context.Context, event *models.Event) (*InsertEventResult, error) {
	query := `
		INSERT INTO events (time, user_id, event_type, source, data, metadata, confidence)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (time, user_id, event_type)
		DO UPDATE SET
			source = EXCLUDED.source,
			data = EXCLUDED.data,
			metadata = EXCLUDED.metadata,
			confidence = EXCLUDED.confidence,
			updated_at = CURRENT_TIMESTAMP
		RETURNING (xmax = 0) AS was_inserted
	`

	var wasInserted bool
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		event.Time,
		event.UserID,
		event.EventType,
		event.Source,
		event.Data,
		event.Metadata,
		event.Confidence,
	).Scan(&wasInserted)

	if err != nil {
		return nil, fmt.Errorf("failed to insert event: %w", err)
	}

	return &InsertEventResult{WasInserted: wasInserted}, nil
}

// GetEventsByUserAndType retrieves events for a user filtered by event type
func (r *EventRepository) GetEventsByUserAndType(
	ctx context.Context,
	userID string,
	eventType string,
	startTime time.Time,
	endTime time.Time,
) ([]models.Event, error) {
	query := `
		SELECT time, user_id, event_type, source, data, metadata, confidence
		FROM events
		WHERE user_id = $1
			AND event_type = $2
			AND time >= $3
			AND time <= $4
		ORDER BY time DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, eventType, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.Time,
			&event.UserID,
			&event.EventType,
			&event.Source,
			&event.Data,
			&event.Metadata,
			&event.Confidence,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// GetEventsByUser retrieves all events for a user within a time range
func (r *EventRepository) GetEventsByUser(
	ctx context.Context,
	userID string,
	startTime time.Time,
	endTime time.Time,
) ([]models.Event, error) {
	query := `
		SELECT time, user_id, event_type, source, data, metadata, confidence
		FROM events
		WHERE user_id = $1
			AND time >= $2
			AND time <= $3
		ORDER BY time DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(
			&event.Time,
			&event.UserID,
			&event.EventType,
			&event.Source,
			&event.Data,
			&event.Metadata,
			&event.Confidence,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// DeleteEvent deletes an event by time, user_id, and event_type
func (r *EventRepository) DeleteEvent(
	ctx context.Context,
	userID string,
	eventType string,
	eventTime time.Time,
) error {
	query := `
		DELETE FROM events
		WHERE user_id = $1
			AND event_type = $2
			AND time = $3
	`

	result, err := r.db.Pool.Exec(ctx, query, userID, eventType, eventTime)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// CountEventsByType counts events by type for a user within a time range
func (r *EventRepository) CountEventsByType(
	ctx context.Context,
	userID string,
	startTime time.Time,
	endTime time.Time,
) (map[string]int64, error) {
	query := `
		SELECT event_type, COUNT(*) as count
		FROM events
		WHERE user_id = $1
			AND time >= $2
			AND time <= $3
		GROUP BY event_type
		ORDER BY count DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query event counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var eventType string
		var count int64
		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan event count: %w", err)
		}
		counts[eventType] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event counts: %w", err)
	}

	return counts, nil
}

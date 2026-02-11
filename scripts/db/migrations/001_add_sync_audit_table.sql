-- Migration: Add sync_audit table for tracking data synchronization
-- This table tracks each sync operation from external sources (Garmin, etc.)

CREATE TABLE IF NOT EXISTS sync_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sync_started_at TIMESTAMPTZ NOT NULL,
    sync_completed_at TIMESTAMPTZ,
    sync_duration_seconds INT,
    user_id UUID NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    target_date VARCHAR(20) NOT NULL,
    records_fetched INT NOT NULL DEFAULT 0,
    records_inserted INT NOT NULL DEFAULT 0,
    records_updated INT NOT NULL DEFAULT 0,
    earliest_timestamp TIMESTAMPTZ,
    latest_timestamp TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failed', 'partial')),
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_sync_audit_user_started ON sync_audit (user_id, sync_started_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_audit_data_type ON sync_audit (data_type, sync_started_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_audit_status ON sync_audit (status, sync_started_at DESC);

-- Grant permissions
GRANT ALL PRIVILEGES ON sync_audit TO healthuser;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'sync_audit table created successfully';
END $$;

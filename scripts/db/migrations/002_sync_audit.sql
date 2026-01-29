-- Migration: Add sync_audit table for tracking data ingestion runs
-- This table provides observability into the sync process

CREATE TABLE IF NOT EXISTS sync_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sync_started_at TIMESTAMPTZ NOT NULL,
    sync_completed_at TIMESTAMPTZ,
    sync_duration_seconds INT,
    user_id UUID NOT NULL,
    data_type VARCHAR(50) NOT NULL, -- 'sleep', 'activity', 'hrv', 'stress'
    target_date DATE NOT NULL,
    records_fetched INT DEFAULT 0,
    records_inserted INT DEFAULT 0,
    records_updated INT DEFAULT 0,
    earliest_timestamp TIMESTAMPTZ,
    latest_timestamp TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL CHECK (status IN ('running', 'success', 'partial', 'failed')),
    error_message TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for querying audit logs
CREATE INDEX IF NOT EXISTS idx_sync_audit_user_date ON sync_audit (user_id, target_date DESC);
CREATE INDEX IF NOT EXISTS idx_sync_audit_status ON sync_audit (status, sync_started_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_audit_data_type ON sync_audit (data_type, sync_started_at DESC);
CREATE INDEX IF NOT EXISTS idx_sync_audit_started_at ON sync_audit (sync_started_at DESC);

-- Grant permissions
GRANT ALL PRIVILEGES ON sync_audit TO healthuser;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Sync audit table created successfully';
END $$;

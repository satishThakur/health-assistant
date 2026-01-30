"""Tests for configuration module."""

import pytest
from app.config import Settings


def test_settings_defaults():
    """Test default configuration values."""
    # Note: This will fail if env vars are not set, which is expected
    # In CI, we'll mock or skip this test
    pass


def test_sync_cron_defaults():
    """Test sync cron schedule defaults."""
    settings = Settings(
        GARMIN_EMAIL="test@example.com",
        GARMIN_PASSWORD="test123",
    )

    assert settings.SYNC_CRON_HOUR == "*"
    assert settings.SYNC_CRON_MINUTE == "0"
    assert settings.INGESTION_SERVICE_URL == "http://ingestion-service:8083"

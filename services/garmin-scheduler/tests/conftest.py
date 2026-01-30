"""Pytest configuration and fixtures."""

import pytest
import os


@pytest.fixture
def mock_env(monkeypatch):
    """Mock environment variables for testing."""
    monkeypatch.setenv("GARMIN_EMAIL", "test@example.com")
    monkeypatch.setenv("GARMIN_PASSWORD", "testpass123")
    monkeypatch.setenv("DEFAULT_USER_ID", "00000000-0000-0000-0000-000000000001")
    monkeypatch.setenv("INGESTION_SERVICE_URL", "http://localhost:8083")


@pytest.fixture
def sample_sleep_data():
    """Sample sleep data for testing."""
    return {
        "sleep_time_seconds": 28800,
        "deep_sleep_seconds": 7200,
        "light_sleep_seconds": 14400,
        "rem_sleep_seconds": 7200,
        "awake_seconds": 0,
        "sleep_scores": {"overall_score": 85},
        "average_hrv": 65.5,
        "sleep_end_timestamp_gmt": "2026-01-28T07:30:00Z"
    }


@pytest.fixture
def sample_activity_data():
    """Sample activity data for testing."""
    return {
        "activity_type": "running",
        "start_time_gmt": "2026-01-28T10:00:00Z",
        "duration_seconds": 2700,
        "distance_meters": 5000,
        "calories": 285,
        "average_heart_rate": 132,
        "max_heart_rate": 168
    }

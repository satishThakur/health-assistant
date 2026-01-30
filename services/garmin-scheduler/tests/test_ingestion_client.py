"""Tests for ingestion client."""

import pytest
from datetime import date
from app.ingestion_client import IngestionClient


def test_ingestion_client_initialization():
    """Test ingestion client initialization."""
    client = IngestionClient(base_url="http://localhost:8083")

    assert client.base_url == "http://localhost:8083"
    assert client.timeout == 30


def test_ingestion_client_custom_timeout():
    """Test ingestion client with custom timeout."""
    client = IngestionClient(
        base_url="http://localhost:8083",
        timeout=60
    )

    assert client.timeout == 60


@pytest.mark.asyncio
async def test_close_client():
    """Test closing the HTTP client."""
    client = IngestionClient(base_url="http://localhost:8083")
    await client.close()
    # Should not raise an error


def test_client_methods_exist():
    """Test that all required methods exist."""
    client = IngestionClient(base_url="http://localhost:8083")

    assert hasattr(client, 'post_sleep_data')
    assert hasattr(client, 'post_activity_data')
    assert hasattr(client, 'post_hrv_data')
    assert hasattr(client, 'post_stress_data')
    assert hasattr(client, 'post_sync_audit')
    assert hasattr(client, 'check_health')

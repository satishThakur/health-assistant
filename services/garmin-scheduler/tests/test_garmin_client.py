"""Tests for Garmin client wrapper."""

import pytest
from datetime import date
from app.garmin_client import GarminClientWrapper


def test_garmin_client_initialization():
    """Test Garmin client can be initialized."""
    client = GarminClientWrapper(
        email="test@example.com",
        password="testpass"
    )

    assert client.email == "test@example.com"
    assert client.password == "testpass"
    assert client.client is None


def test_connect_requires_credentials():
    """Test that connect requires valid credentials."""
    client = GarminClientWrapper(
        email="test@example.com",
        password="testpass"
    )

    # Connection will fail with invalid credentials
    # This test just ensures the method exists
    assert hasattr(client, 'connect')
    assert hasattr(client, 'get_sleep_data')
    assert hasattr(client, 'get_activity_data')
    assert hasattr(client, 'get_hrv_data')
    assert hasattr(client, 'get_stress_data')


def test_data_methods_require_connection():
    """Test that data methods check for connection."""
    client = GarminClientWrapper(
        email="test@example.com",
        password="testpass"
    )

    target_date = date(2026, 1, 28)

    # Should raise RuntimeError when not connected
    with pytest.raises(RuntimeError, match="Client not connected"):
        client.get_sleep_data(target_date)

    with pytest.raises(RuntimeError, match="Client not connected"):
        client.get_activity_data(target_date)

    with pytest.raises(RuntimeError, match="Client not connected"):
        client.get_hrv_data(target_date)

    with pytest.raises(RuntimeError, match="Client not connected"):
        client.get_stress_data(target_date)

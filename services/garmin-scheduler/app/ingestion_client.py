"""HTTP client for posting data to the Go ingestion service."""

import logging
from typing import Dict, Any, Optional
from datetime import date

import httpx

logger = logging.getLogger(__name__)


class IngestionClient:
    """Client for sending data to the Go ingestion service."""

    def __init__(self, base_url: str, timeout: int = 30):
        """
        Initialize ingestion client.

        Args:
            base_url: Base URL of the ingestion service (e.g., http://ingestion-service:8083)
            timeout: Request timeout in seconds
        """
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self.client = httpx.AsyncClient(timeout=timeout)

    async def close(self):
        """Close the HTTP client."""
        await self.client.aclose()

    async def post_sleep_data(
        self,
        user_id: str,
        target_date: date,
        sleep_data: Dict[str, Any],
    ) -> bool:
        """
        Post sleep data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the sleep data
            sleep_data: Sleep metrics dictionary

        Returns:
            True if successful, False otherwise
        """
        url = f"{self.base_url}/api/v1/garmin/ingest/sleep"
        payload = {
            "user_id": user_id,
            "date": target_date.isoformat(),
            "sleep_data": sleep_data,
        }

        try:
            logger.info(f"Posting sleep data for user {user_id} on {target_date}")
            response = await self.client.post(url, json=payload)
            response.raise_for_status()
            logger.info(f"Successfully posted sleep data: {response.status_code}")
            return True
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting sleep data: {e.response.status_code} - {e.response.text}")
            return False
        except Exception as e:
            logger.error(f"Error posting sleep data: {e}")
            return False

    async def post_activity_data(
        self,
        user_id: str,
        target_date: date,
        activity_data: Dict[str, Any],
    ) -> bool:
        """
        Post activity data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the activity
            activity_data: Activity metrics dictionary

        Returns:
            True if successful, False otherwise
        """
        url = f"{self.base_url}/api/v1/garmin/ingest/activity"
        payload = {
            "user_id": user_id,
            "date": target_date.isoformat(),
            "activity_data": activity_data,
        }

        try:
            logger.info(f"Posting activity data for user {user_id} on {target_date}")
            response = await self.client.post(url, json=payload)
            response.raise_for_status()
            logger.info(f"Successfully posted activity data: {response.status_code}")
            return True
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting activity data: {e.response.status_code} - {e.response.text}")
            return False
        except Exception as e:
            logger.error(f"Error posting activity data: {e}")
            return False

    async def post_hrv_data(
        self,
        user_id: str,
        target_date: date,
        hrv_data: Dict[str, Any],
    ) -> bool:
        """
        Post HRV data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the HRV data
            hrv_data: HRV metrics dictionary

        Returns:
            True if successful, False otherwise
        """
        url = f"{self.base_url}/api/v1/garmin/ingest/hrv"
        payload = {
            "user_id": user_id,
            "date": target_date.isoformat(),
            "hrv_data": hrv_data,
        }

        try:
            logger.info(f"Posting HRV data for user {user_id} on {target_date}")
            response = await self.client.post(url, json=payload)
            response.raise_for_status()
            logger.info(f"Successfully posted HRV data: {response.status_code}")
            return True
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting HRV data: {e.response.status_code} - {e.response.text}")
            return False
        except Exception as e:
            logger.error(f"Error posting HRV data: {e}")
            return False

    async def post_stress_data(
        self,
        user_id: str,
        target_date: date,
        stress_data: Dict[str, Any],
    ) -> bool:
        """
        Post stress data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the stress data
            stress_data: Stress metrics dictionary

        Returns:
            True if successful, False otherwise
        """
        url = f"{self.base_url}/api/v1/garmin/ingest/stress"
        payload = {
            "user_id": user_id,
            "date": target_date.isoformat(),
            "stress_data": stress_data,
        }

        try:
            logger.info(f"Posting stress data for user {user_id} on {target_date}")
            response = await self.client.post(url, json=payload)
            response.raise_for_status()
            logger.info(f"Successfully posted stress data: {response.status_code}")
            return True
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting stress data: {e.response.status_code} - {e.response.text}")
            return False
        except Exception as e:
            logger.error(f"Error posting stress data: {e}")
            return False

    async def check_health(self) -> bool:
        """
        Check if the ingestion service is healthy.

        Returns:
            True if service is healthy, False otherwise
        """
        url = f"{self.base_url}/health"

        try:
            response = await self.client.get(url)
            response.raise_for_status()
            return True
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return False

"""HTTP client for posting data to the Go ingestion service."""

import logging
from typing import Dict, Any, Optional
from datetime import date, datetime

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
    ) -> Optional[Dict[str, Any]]:
        """
        Post sleep data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the sleep data
            sleep_data: Sleep metrics dictionary

        Returns:
            Response dict with status and was_inserted, or None if failed
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
            result = response.json()
            logger.info(f"Successfully posted sleep data: {result.get('action', 'unknown')}")
            return result
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting sleep data: {e.response.status_code} - {e.response.text}")
            return None
        except Exception as e:
            logger.error(f"Error posting sleep data: {e}")
            return None

    async def post_activity_data(
        self,
        user_id: str,
        target_date: date,
        activity_data: Dict[str, Any],
    ) -> Optional[Dict[str, Any]]:
        """
        Post activity data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the activity
            activity_data: Activity metrics dictionary

        Returns:
            Response dict with status and was_inserted, or None if failed
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
            result = response.json()
            logger.info(f"Successfully posted activity data: {result.get('action', 'unknown')}")
            return result
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting activity data: {e.response.status_code} - {e.response.text}")
            return None
        except Exception as e:
            logger.error(f"Error posting activity data: {e}")
            return None

    async def post_hrv_data(
        self,
        user_id: str,
        target_date: date,
        hrv_data: Dict[str, Any],
    ) -> Optional[Dict[str, Any]]:
        """
        Post HRV data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the HRV data
            hrv_data: HRV metrics dictionary

        Returns:
            Response dict with status and was_inserted, or None if failed
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
            result = response.json()
            logger.info(f"Successfully posted HRV data: {result.get('action', 'unknown')}")
            return result
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting HRV data: {e.response.status_code} - {e.response.text}")
            return None
        except Exception as e:
            logger.error(f"Error posting HRV data: {e}")
            return None

    async def post_stress_data(
        self,
        user_id: str,
        target_date: date,
        stress_data: Dict[str, Any],
    ) -> Optional[Dict[str, Any]]:
        """
        Post stress data to ingestion service.

        Args:
            user_id: User UUID
            target_date: Date of the stress data
            stress_data: Stress metrics dictionary

        Returns:
            Response dict with status and was_inserted, or None if failed
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
            result = response.json()
            logger.info(f"Successfully posted stress data: {result.get('action', 'unknown')}")
            return result
        except httpx.HTTPStatusError as e:
            logger.error(f"HTTP error posting stress data: {e.response.status_code} - {e.response.text}")
            return None
        except Exception as e:
            logger.error(f"Error posting stress data: {e}")
            return None

    async def post_sync_audit(
        self,
        user_id: str,
        data_type: str,
        target_date: date,
        sync_started_at: datetime,
        sync_completed_at: datetime,
        sync_duration_seconds: int,
        records_fetched: int,
        records_inserted: int,
        records_updated: int,
        earliest_timestamp: Optional[str],
        latest_timestamp: Optional[str],
        status: str,
        error_message: Optional[str] = None,
    ) -> bool:
        """
        Post sync audit record to ingestion service.

        Args:
            user_id: User UUID
            data_type: Type of data synced
            target_date: Date of the data
            sync_started_at: When sync started
            sync_completed_at: When sync completed
            sync_duration_seconds: Duration in seconds
            records_fetched: Number of records fetched from Garmin
            records_inserted: Number of new records inserted
            records_updated: Number of existing records updated
            earliest_timestamp: Earliest timestamp in the data
            latest_timestamp: Latest timestamp in the data
            status: Sync status ('success', 'partial', 'failed')
            error_message: Error message if failed

        Returns:
            True if successful, False otherwise
        """
        url = f"{self.base_url}/api/v1/audit/sync"
        payload = {
            "sync_started_at": sync_started_at.isoformat() + "Z",
            "sync_completed_at": sync_completed_at.isoformat() + "Z",
            "sync_duration_seconds": sync_duration_seconds,
            "user_id": user_id,
            "data_type": data_type,
            "target_date": target_date.isoformat(),
            "records_fetched": records_fetched,
            "records_inserted": records_inserted,
            "records_updated": records_updated,
            "status": status,
        }

        if earliest_timestamp:
            payload["earliest_timestamp"] = earliest_timestamp
        if latest_timestamp:
            payload["latest_timestamp"] = latest_timestamp
        if error_message:
            payload["error_message"] = error_message

        try:
            response = await self.client.post(url, json=payload)
            response.raise_for_status()
            logger.debug(f"Audit logged for {data_type} on {target_date}")
            return True
        except Exception as e:
            logger.error(f"Failed to post sync audit: {e}")
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

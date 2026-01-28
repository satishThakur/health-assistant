"""APScheduler-based Garmin data sync scheduler."""

import asyncio
import logging
from datetime import date, timedelta
from typing import Optional

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.cron import CronTrigger

from .config import settings
from .garmin_client import GarminClientWrapper
from .ingestion_client import IngestionClient

logger = logging.getLogger(__name__)


class GarminSyncScheduler:
    """Scheduler for periodic Garmin data synchronization."""

    def __init__(self):
        """Initialize the scheduler."""
        self.scheduler = AsyncIOScheduler()
        self.garmin_client: Optional[GarminClientWrapper] = None
        self.ingestion_client: Optional[IngestionClient] = None
        self._sync_lock = asyncio.Lock()

    def start(self):
        """Start the scheduler."""
        # Initialize clients
        self.garmin_client = GarminClientWrapper(
            email=settings.GARMIN_EMAIL,
            password=settings.GARMIN_PASSWORD,
        )
        self.ingestion_client = IngestionClient(base_url=settings.INGESTION_SERVICE_URL)

        # Add sync job with cron schedule
        self.scheduler.add_job(
            self.sync_garmin_data,
            CronTrigger(
                hour=settings.SYNC_CRON_HOUR,
                minute=settings.SYNC_CRON_MINUTE,
            ),
            id="garmin_sync",
            name="Sync Garmin data to ingestion service",
            replace_existing=True,
        )

        logger.info(
            f"Scheduler configured: sync at {settings.SYNC_CRON_HOUR}:{settings.SYNC_CRON_MINUTE}"
        )

        # Start the scheduler
        self.scheduler.start()
        logger.info("Scheduler started")

    async def stop(self):
        """Stop the scheduler and cleanup."""
        logger.info("Stopping scheduler...")
        self.scheduler.shutdown(wait=True)

        if self.ingestion_client:
            await self.ingestion_client.close()

        logger.info("Scheduler stopped")

    async def sync_garmin_data(self):
        """
        Sync Garmin data to ingestion service.

        This is the main scheduled job that:
        1. Connects to Garmin
        2. Fetches data for yesterday and today
        3. Posts to Go ingestion service
        """
        # Prevent concurrent sync runs
        if self._sync_lock.locked():
            logger.warning("Sync already in progress, skipping this run")
            return

        async with self._sync_lock:
            logger.info("Starting Garmin data sync")

            try:
                # Check if ingestion service is healthy
                if not await self.ingestion_client.check_health():
                    logger.error("Ingestion service is not healthy, aborting sync")
                    return

                # Connect to Garmin
                self.garmin_client.connect()

                # Sync data for yesterday and today
                today = date.today()
                yesterday = today - timedelta(days=1)

                for target_date in [yesterday, today]:
                    await self._sync_date(target_date)

                logger.info("Garmin data sync completed successfully")

            except Exception as e:
                logger.error(f"Error during Garmin sync: {e}", exc_info=True)

    async def _sync_date(self, target_date: date):
        """
        Sync all data types for a specific date.

        Args:
            target_date: Date to sync data for
        """
        logger.info(f"Syncing data for {target_date}")
        user_id = settings.DEFAULT_USER_ID

        # Sync sleep data
        try:
            sleep_data = self.garmin_client.get_sleep_data(target_date)
            if sleep_data:
                success = await self.ingestion_client.post_sleep_data(
                    user_id=user_id,
                    target_date=target_date,
                    sleep_data=sleep_data,
                )
                if success:
                    logger.info(f"Sleep data synced for {target_date}")
        except Exception as e:
            logger.error(f"Error syncing sleep data for {target_date}: {e}")

        # Sync activity data
        try:
            activities = self.garmin_client.get_activity_data(target_date)
            if activities:
                for activity in activities:
                    success = await self.ingestion_client.post_activity_data(
                        user_id=user_id,
                        target_date=target_date,
                        activity_data=activity,
                    )
                    if success:
                        logger.info(f"Activity synced for {target_date}: {activity['activity_type']}")
        except Exception as e:
            logger.error(f"Error syncing activity data for {target_date}: {e}")

        # Sync HRV data
        try:
            hrv_data = self.garmin_client.get_hrv_data(target_date)
            if hrv_data:
                success = await self.ingestion_client.post_hrv_data(
                    user_id=user_id,
                    target_date=target_date,
                    hrv_data=hrv_data,
                )
                if success:
                    logger.info(f"HRV data synced for {target_date}")
        except Exception as e:
            logger.error(f"Error syncing HRV data for {target_date}: {e}")

        # Sync stress data
        try:
            stress_data = self.garmin_client.get_stress_data(target_date)
            if stress_data:
                success = await self.ingestion_client.post_stress_data(
                    user_id=user_id,
                    target_date=target_date,
                    stress_data=stress_data,
                )
                if success:
                    logger.info(f"Stress data synced for {target_date}")
        except Exception as e:
            logger.error(f"Error syncing stress data for {target_date}: {e}")

    async def trigger_manual_sync(self):
        """Manually trigger a sync (for testing/on-demand sync)."""
        logger.info("Manual sync triggered")
        await self.sync_garmin_data()

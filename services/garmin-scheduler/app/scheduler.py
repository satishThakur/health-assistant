"""APScheduler-based Garmin data sync scheduler."""

import asyncio
import logging
from datetime import date, datetime, timedelta
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
        Sync all data types for a specific date with audit tracking.

        Args:
            target_date: Date to sync data for
        """
        logger.info(f"Syncing data for {target_date}")
        user_id = settings.DEFAULT_USER_ID

        # Sync sleep data with audit
        await self._sync_data_type(
            data_type="sleep",
            target_date=target_date,
            user_id=user_id,
        )

        # Sync activity data with audit
        await self._sync_data_type(
            data_type="activity",
            target_date=target_date,
            user_id=user_id,
        )

        # Sync HRV data with audit
        await self._sync_data_type(
            data_type="hrv",
            target_date=target_date,
            user_id=user_id,
        )

        # Sync stress data with audit
        await self._sync_data_type(
            data_type="stress",
            target_date=target_date,
            user_id=user_id,
        )

        # Sync daily stats with audit
        await self._sync_data_type(
            data_type="daily_stats",
            target_date=target_date,
            user_id=user_id,
        )

        # Sync body battery with audit
        await self._sync_data_type(
            data_type="body_battery",
            target_date=target_date,
            user_id=user_id,
        )

    async def _sync_data_type(self, data_type: str, target_date: date, user_id: str):
        """
        Sync a specific data type with full audit tracking.

        Args:
            data_type: Type of data ('sleep', 'activity', 'hrv', 'stress')
            target_date: Date to sync
            user_id: User ID
        """
        sync_started_at = datetime.utcnow()
        records_fetched = 0
        records_inserted = 0
        records_updated = 0
        status = "success"
        error_message = None
        earliest_timestamp = None
        latest_timestamp = None

        try:
            # Fetch data from Garmin
            if data_type == "sleep":
                data = self.garmin_client.get_sleep_data(target_date)
                if data:
                    records_fetched = 1
                    response = await self.ingestion_client.post_sleep_data(
                        user_id=user_id,
                        target_date=target_date,
                        sleep_data=data,
                    )
                    if response and response.get("was_inserted"):
                        records_inserted = 1
                    else:
                        records_updated = 1

                    # Extract timestamp if available
                    if "sleep_end_timestamp_gmt" in data:
                        earliest_timestamp = latest_timestamp = data["sleep_end_timestamp_gmt"]

            elif data_type == "activity":
                activities = self.garmin_client.get_activity_data(target_date)
                if activities:
                    records_fetched = len(activities)
                    for activity in activities:
                        response = await self.ingestion_client.post_activity_data(
                            user_id=user_id,
                            target_date=target_date,
                            activity_data=activity,
                        )
                        if response:
                            if response.get("was_inserted"):
                                records_inserted += 1
                            else:
                                records_updated += 1

                            # Track timestamps
                            if "start_time_gmt" in activity:
                                ts = activity["start_time_gmt"]
                                if earliest_timestamp is None or ts < earliest_timestamp:
                                    earliest_timestamp = ts
                                if latest_timestamp is None or ts > latest_timestamp:
                                    latest_timestamp = ts

            elif data_type == "hrv":
                data = self.garmin_client.get_hrv_data(target_date)
                if data:
                    records_fetched = 1
                    response = await self.ingestion_client.post_hrv_data(
                        user_id=user_id,
                        target_date=target_date,
                        hrv_data=data,
                    )
                    if response and response.get("was_inserted"):
                        records_inserted = 1
                    else:
                        records_updated = 1

            elif data_type == "stress":
                data = self.garmin_client.get_stress_data(target_date)
                if data:
                    records_fetched = 1
                    response = await self.ingestion_client.post_stress_data(
                        user_id=user_id,
                        target_date=target_date,
                        stress_data=data,
                    )
                    if response and response.get("was_inserted"):
                        records_inserted = 1
                    else:
                        records_updated = 1

            elif data_type == "daily_stats":
                data = self.garmin_client.get_daily_stats(target_date)
                if data:
                    records_fetched = 1
                    response = await self.ingestion_client.post_daily_stats(
                        user_id=user_id,
                        target_date=target_date,
                        daily_stats_data=data,
                    )
                    if response and response.get("was_inserted"):
                        records_inserted = 1
                    else:
                        records_updated = 1

            elif data_type == "body_battery":
                data = self.garmin_client.get_body_battery(target_date)
                if data:
                    records_fetched = 1
                    response = await self.ingestion_client.post_body_battery(
                        user_id=user_id,
                        target_date=target_date,
                        body_battery_data=data,
                    )
                    if response and response.get("was_inserted"):
                        records_inserted = 1
                    else:
                        records_updated = 1

            if records_fetched > 0:
                logger.info(
                    f"{data_type.capitalize()} sync for {target_date}: "
                    f"fetched={records_fetched}, inserted={records_inserted}, updated={records_updated}"
                )

        except Exception as e:
            status = "failed"
            error_message = str(e)
            logger.error(f"Error syncing {data_type} data for {target_date}: {e}")

        # Record audit
        sync_completed_at = datetime.utcnow()
        sync_duration_seconds = int((sync_completed_at - sync_started_at).total_seconds())

        await self.ingestion_client.post_sync_audit(
            user_id=user_id,
            data_type=data_type,
            target_date=target_date,
            sync_started_at=sync_started_at,
            sync_completed_at=sync_completed_at,
            sync_duration_seconds=sync_duration_seconds,
            records_fetched=records_fetched,
            records_inserted=records_inserted,
            records_updated=records_updated,
            earliest_timestamp=earliest_timestamp,
            latest_timestamp=latest_timestamp,
            status=status,
            error_message=error_message,
        )

    async def trigger_manual_sync(self):
        """Manually trigger a sync (for testing/on-demand sync)."""
        logger.info("Manual sync triggered")
        await self.sync_garmin_data()

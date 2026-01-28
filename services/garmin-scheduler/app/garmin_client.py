"""Garmin API client wrapper for fetching health data."""

import logging
from datetime import date, datetime
from typing import Optional, Dict, Any

from garminconnect import Garmin, GarminConnectAuthenticationError, GarminConnectConnectionError

logger = logging.getLogger(__name__)


class GarminClientWrapper:
    """Wrapper around garminconnect library for fetching health data."""

    def __init__(self, email: str, password: str):
        """Initialize Garmin client with credentials."""
        self.email = email
        self.password = password
        self.client: Optional[Garmin] = None

    def connect(self) -> None:
        """Authenticate with Garmin Connect."""
        try:
            logger.info("Connecting to Garmin Connect...")
            self.client = Garmin(self.email, self.password)
            self.client.login()
            logger.info("Successfully connected to Garmin Connect")
        except GarminConnectAuthenticationError:
            logger.error("Authentication failed - check credentials")
            raise
        except GarminConnectConnectionError as e:
            logger.error(f"Connection error: {e}")
            raise

    def get_sleep_data(self, target_date: date) -> Optional[Dict[str, Any]]:
        """
        Fetch sleep data for a specific date and transform to ingestion format.

        Returns:
            Dict with keys: sleep_time_seconds, deep_sleep_seconds, light_sleep_seconds,
            rem_sleep_seconds, awake_seconds, sleep_scores, average_hrv, sleep_end_timestamp_gmt
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching sleep data for {date_str}")

            sleep_data = self.client.get_sleep_data(date_str)
            if not sleep_data:
                logger.info(f"No sleep data available for {date_str}")
                return None

            # Extract sleep metrics
            daily_sleep = sleep_data.get("dailySleepDTO", {})
            sleep_levels = daily_sleep.get("sleepLevels", {})

            # Transform to ingestion format
            transformed = {
                "sleep_time_seconds": daily_sleep.get("sleepTimeSeconds", 0),
                "deep_sleep_seconds": daily_sleep.get("deepSleepSeconds", 0),
                "light_sleep_seconds": daily_sleep.get("lightSleepSeconds", 0),
                "rem_sleep_seconds": daily_sleep.get("remSleepSeconds", 0),
                "awake_seconds": daily_sleep.get("awakeSleepSeconds", 0),
                "sleep_scores": {
                    "overall_score": daily_sleep.get("sleepScores", {}).get("overall", {}).get("value", 0)
                },
            }

            # Add HRV if available
            if "averageHRV" in daily_sleep:
                transformed["average_hrv"] = daily_sleep.get("averageHRV")

            # Add sleep end timestamp
            if "sleepEndTimestampGMT" in daily_sleep:
                # Convert Unix timestamp (milliseconds) to ISO format
                timestamp_ms = daily_sleep["sleepEndTimestampGMT"]
                dt = datetime.fromtimestamp(timestamp_ms / 1000)
                transformed["sleep_end_timestamp_gmt"] = dt.isoformat() + "Z"

            logger.info(f"Successfully fetched sleep data for {date_str}: {transformed['sleep_time_seconds']/60:.0f} min")
            return transformed

        except Exception as e:
            logger.error(f"Error fetching sleep data for {target_date}: {e}")
            return None

    def get_activity_data(self, target_date: date) -> Optional[list[Dict[str, Any]]]:
        """
        Fetch activity data for a specific date and transform to ingestion format.

        Returns:
            List of activities with keys: activity_type, start_time_gmt, duration_seconds,
            distance_meters, calories, average_heart_rate, max_heart_rate
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching activity data for {date_str}")

            # Get activities for the date
            activities = self.client.get_activities_by_date(date_str, date_str)
            if not activities:
                logger.info(f"No activities available for {date_str}")
                return None

            transformed_activities = []
            for activity in activities:
                # Extract activity metrics
                activity_data = {
                    "activity_type": activity.get("activityType", {}).get("typeKey", "unknown"),
                    "duration_seconds": activity.get("duration", 0),
                    "distance_meters": activity.get("distance", 0),
                    "calories": activity.get("calories", 0),
                }

                # Add HR data if available
                if "averageHR" in activity:
                    activity_data["average_heart_rate"] = activity.get("averageHR")
                if "maxHR" in activity:
                    activity_data["max_heart_rate"] = activity.get("maxHR")

                # Add start time
                if "startTimeGMT" in activity:
                    # Convert Unix timestamp (milliseconds) to ISO format
                    timestamp_ms = activity["startTimeGMT"]
                    dt = datetime.fromtimestamp(timestamp_ms / 1000)
                    activity_data["start_time_gmt"] = dt.isoformat() + "Z"

                transformed_activities.append(activity_data)

            logger.info(f"Successfully fetched {len(transformed_activities)} activities for {date_str}")
            return transformed_activities

        except Exception as e:
            logger.error(f"Error fetching activity data for {target_date}: {e}")
            return None

    def get_hrv_data(self, target_date: date) -> Optional[Dict[str, Any]]:
        """
        Fetch HRV data for a specific date and transform to ingestion format.

        Returns:
            Dict with keys: average_hrv, max_hrv, min_hrv (if available)
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching HRV data for {date_str}")

            # HRV is typically part of sleep data or stress data
            # Try to get from sleep data first
            sleep_data = self.client.get_sleep_data(date_str)
            if sleep_data:
                daily_sleep = sleep_data.get("dailySleepDTO", {})
                if "averageHRV" in daily_sleep:
                    transformed = {
                        "average_hrv": daily_sleep.get("averageHRV"),
                    }
                    logger.info(f"Successfully fetched HRV data for {date_str}: {transformed['average_hrv']}")
                    return transformed

            # Try stress data as alternative source
            try:
                stress_data = self.client.get_stress_data(date_str)
                if stress_data and "avgStressLevel" in stress_data:
                    # Some devices provide HRV through stress API
                    logger.info(f"No HRV data in sleep data for {date_str}")
                    return None
            except:
                pass

            logger.info(f"No HRV data available for {date_str}")
            return None

        except Exception as e:
            logger.error(f"Error fetching HRV data for {target_date}: {e}")
            return None

    def get_stress_data(self, target_date: date) -> Optional[Dict[str, Any]]:
        """
        Fetch stress data for a specific date and transform to ingestion format.

        Returns:
            Dict with keys: average_stress_level, max_stress_level, rest_stress_duration
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching stress data for {date_str}")

            stress_data = self.client.get_stress_data(date_str)
            if not stress_data:
                logger.info(f"No stress data available for {date_str}")
                return None

            # Transform to ingestion format
            transformed = {}

            if "avgStressLevel" in stress_data:
                transformed["average_stress_level"] = stress_data.get("avgStressLevel")
            if "maxStressLevel" in stress_data:
                transformed["max_stress_level"] = stress_data.get("maxStressLevel")
            if "restStressDuration" in stress_data:
                transformed["rest_stress_duration"] = stress_data.get("restStressDuration")

            if transformed:
                logger.info(f"Successfully fetched stress data for {date_str}")
                return transformed

            logger.info(f"No stress data available for {date_str}")
            return None

        except Exception as e:
            logger.error(f"Error fetching stress data for {target_date}: {e}")
            return None

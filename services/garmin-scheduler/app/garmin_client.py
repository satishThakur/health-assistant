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
            logger.info(f"Email: {self.email}")
            logger.info(f"Password length: {len(self.password)} chars")
            logger.info(f"Password starts with: {self.password[:2] if len(self.password) >= 2 else 'N/A'}")

            self.client = Garmin(self.email, self.password)
            logger.info("Garmin client created, attempting login...")

            self.client.login()

            logger.info("Successfully connected to Garmin Connect")
            logger.info(f"Display name: {self.client.display_name if hasattr(self.client, 'display_name') else 'N/A'}")

        except GarminConnectAuthenticationError as e:
            logger.error(f"Authentication failed - check credentials: {e}")
            raise
        except GarminConnectConnectionError as e:
            logger.error(f"Connection error: {e}")
            raise
        except Exception as e:
            logger.error(f"Unexpected error during login: {type(e).__name__}: {e}")
            logger.error(f"Full error: {repr(e)}")
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
            Dict with keys: last_night_avg, weekly_avg, status
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching HRV data for {date_str}")

            # Use dedicated HRV method
            hrv_data = self.client.get_hrv_data(date_str)
            if hrv_data:
                transformed = {}

                if "lastNightAvg" in hrv_data:
                    transformed["last_night_avg"] = hrv_data.get("lastNightAvg")
                if "weeklyAvg" in hrv_data:
                    transformed["weekly_avg"] = hrv_data.get("weeklyAvg")
                if "hrvStatus" in hrv_data:
                    transformed["status"] = hrv_data.get("hrvStatus")

                if transformed:
                    logger.info(f"Successfully fetched HRV data for {date_str}")
                    return transformed

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

    def get_daily_stats(self, target_date: date) -> Optional[Dict[str, Any]]:
        """
        Fetch daily stats for a specific date and transform to ingestion format.

        Returns:
            Dict with keys: steps, calories, distance_meters, active_minutes,
            resting_heart_rate, max_heart_rate
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching daily stats for {date_str}")

            stats = self.client.get_stats(date_str)
            if not stats:
                logger.info(f"No daily stats available for {date_str}")
                return None

            # Transform to ingestion format
            transformed = {}

            if "totalSteps" in stats:
                transformed["steps"] = stats.get("totalSteps")
            if "totalKilocalories" in stats:
                transformed["calories"] = stats.get("totalKilocalories")
            if "totalDistanceMeters" in stats:
                transformed["distance_meters"] = stats.get("totalDistanceMeters")
            if "activeKilocalories" in stats:
                transformed["active_calories"] = stats.get("activeKilocalories")
            if "bmrKilocalories" in stats:
                transformed["bmr_calories"] = stats.get("bmrKilocalories")
            if "minHeartRate" in stats:
                transformed["min_heart_rate"] = stats.get("minHeartRate")
            if "maxHeartRate" in stats:
                transformed["max_heart_rate"] = stats.get("maxHeartRate")
            if "restingHeartRate" in stats:
                transformed["resting_heart_rate"] = stats.get("restingHeartRate")
            if "moderateIntensityMinutes" in stats:
                transformed["moderate_intensity_minutes"] = stats.get("moderateIntensityMinutes")
            if "vigorousIntensityMinutes" in stats:
                transformed["vigorous_intensity_minutes"] = stats.get("vigorousIntensityMinutes")

            if transformed:
                logger.info(f"Successfully fetched daily stats for {date_str}: {transformed.get('steps', 0)} steps")
                return transformed

            logger.info(f"No daily stats available for {date_str}")
            return None

        except Exception as e:
            logger.error(f"Error fetching daily stats for {target_date}: {e}")
            return None

    def get_body_battery(self, target_date: date) -> Optional[Dict[str, Any]]:
        """
        Fetch body battery data for a specific date and transform to ingestion format.

        Returns:
            Dict with keys: charged, drained, highest_value, lowest_value
        """
        if not self.client:
            raise RuntimeError("Client not connected. Call connect() first.")

        try:
            date_str = target_date.isoformat()
            logger.info(f"Fetching body battery for {date_str}")

            # Body battery API returns data for a date range
            body_battery = self.client.get_body_battery(date_str, date_str)
            if not body_battery or len(body_battery) == 0:
                logger.info(f"No body battery data available for {date_str}")
                return None

            # Get the first day's data
            day_data = body_battery[0] if isinstance(body_battery, list) else body_battery

            # Transform to ingestion format
            transformed = {}

            if "charged" in day_data:
                transformed["charged"] = day_data.get("charged")
            if "drained" in day_data:
                transformed["drained"] = day_data.get("drained")
            if "highest" in day_data:
                transformed["highest_value"] = day_data.get("highest")
            if "lowest" in day_data:
                transformed["lowest_value"] = day_data.get("lowest")

            if transformed:
                logger.info(f"Successfully fetched body battery for {date_str}")
                return transformed

            logger.info(f"No body battery data available for {date_str}")
            return None

        except Exception as e:
            logger.error(f"Error fetching body battery for {target_date}: {e}")
            return None

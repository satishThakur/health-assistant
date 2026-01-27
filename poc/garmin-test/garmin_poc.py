#!/usr/bin/env python3
"""
Simple POC to test python-garminconnect library.

This script logs into Garmin Connect and fetches your personal health data
to validate that the unofficial API works before integrating into the main system.
"""

import sys
import json
from datetime import date, timedelta
from getpass import getpass
from garminconnect import Garmin, GarminConnectAuthenticationError, GarminConnectConnectionError


def print_section(title):
    """Print a formatted section header."""
    print(f"\n{'='*60}")
    print(f"  {title}")
    print(f"{'='*60}\n")


def print_json(data, indent=2):
    """Pretty print JSON data."""
    print(json.dumps(data, indent=indent, default=str))


def main():
    print_section("Garmin Connect POC - Data Fetcher")

    # Get credentials
    print("Enter your Garmin Connect credentials:")
    email = input("Email: ").strip()
    password = getpass("Password: ").strip()

    if not email or not password:
        print("Error: Email and password are required!")
        sys.exit(1)

    try:
        # Initialize Garmin client
        print("\n🔄 Connecting to Garmin Connect...")
        client = Garmin(email, password)

        # Login
        client.login()
        print("✅ Login successful!\n")

        # Get today's date and yesterday
        today = date.today()
        yesterday = today - timedelta(days=1)

        # Fetch user profile
        print_section("User Profile")
        try:
            profile = client.get_full_name()
            print(f"Name: {profile}")
        except Exception as e:
            print(f"Could not fetch profile: {e}")

        # Fetch daily summary stats
        print_section(f"Daily Stats - {today}")
        try:
            stats = client.get_stats(today.isoformat())
            print(f"Steps: {stats.get('totalSteps', 'N/A')}")
            print(f"Calories: {stats.get('totalKilocalories', 'N/A')}")
            print(f"Distance (meters): {stats.get('totalDistanceMeters', 'N/A')}")
            print(f"Active minutes: {stats.get('activeKilocalories', 'N/A')}")
            print(f"Resting HR: {stats.get('restingHeartRate', 'N/A')}")
            print(f"Max HR: {stats.get('maxHeartRate', 'N/A')}")
        except Exception as e:
            print(f"Error fetching stats: {e}")

        # Fetch sleep data (usually available for yesterday)
        print_section(f"Sleep Data - {yesterday}")
        try:
            sleep = client.get_sleep_data(yesterday.isoformat())
            if sleep:
                daily_sleep = sleep.get('dailySleepDTO', {})
                print(f"Sleep start: {daily_sleep.get('sleepStartTimestampLocal', 'N/A')}")
                print(f"Sleep end: {daily_sleep.get('sleepEndTimestampLocal', 'N/A')}")
                print(f"Total sleep time (sec): {daily_sleep.get('sleepTimeSeconds', 'N/A')}")
                print(f"Deep sleep (sec): {daily_sleep.get('deepSleepSeconds', 'N/A')}")
                print(f"Light sleep (sec): {daily_sleep.get('lightSleepSeconds', 'N/A')}")
                print(f"REM sleep (sec): {daily_sleep.get('remSleepSeconds', 'N/A')}")
                print(f"Awake time (sec): {daily_sleep.get('awakeSleepSeconds', 'N/A')}")

                # Sleep levels
                levels = sleep.get('sleepLevels', [])
                if levels:
                    print(f"\nSleep score: {daily_sleep.get('sleepScores', {}).get('overall', {}).get('value', 'N/A')}")
        except Exception as e:
            print(f"Error fetching sleep data: {e}")

        # Fetch recent activities
        print_section(f"Recent Activities (last 5)")
        try:
            activities = client.get_activities(0, 5)  # Get last 5 activities
            if activities:
                for i, activity in enumerate(activities, 1):
                    print(f"\n{i}. {activity.get('activityName', 'Unknown')}")
                    print(f"   Type: {activity.get('activityType', {}).get('typeKey', 'N/A')}")
                    print(f"   Date: {activity.get('startTimeLocal', 'N/A')}")
                    print(f"   Duration (sec): {activity.get('duration', 'N/A')}")
                    print(f"   Distance (meters): {activity.get('distance', 'N/A')}")
                    print(f"   Calories: {activity.get('calories', 'N/A')}")
                    print(f"   Average HR: {activity.get('averageHR', 'N/A')}")
                    print(f"   Max HR: {activity.get('maxHR', 'N/A')}")
            else:
                print("No activities found")
        except Exception as e:
            print(f"Error fetching activities: {e}")

        # Fetch HRV data (if available on your device)
        print_section(f"HRV Data - {yesterday}")
        try:
            hrv = client.get_hrv_data(yesterday.isoformat())
            if hrv:
                print(f"HRV Status: {hrv.get('hrvStatus', 'N/A')}")
                print(f"Last Night Average: {hrv.get('lastNightAvg', 'N/A')}")
                print(f"Weekly Average: {hrv.get('weeklyAvg', 'N/A')}")
                print(f"Baseline: {hrv.get('baseline', {}).get('lowUpper', 'N/A')} - {hrv.get('baseline', {}).get('balancedHigh', 'N/A')}")
            else:
                print("No HRV data available")
        except Exception as e:
            print(f"Error fetching HRV data: {e}")
            print("Note: HRV requires compatible Garmin device (e.g., Fenix, Forerunner 945+)")

        # Fetch body battery (if available)
        print_section(f"Body Battery - {today}")
        try:
            body_battery = client.get_body_battery(today.isoformat(), today.isoformat())
            if body_battery:
                # Body battery returns list of readings throughout the day
                readings = body_battery[0] if isinstance(body_battery, list) and len(body_battery) > 0 else []
                if readings:
                    charged = readings.get('charged', 'N/A')
                    drained = readings.get('drained', 'N/A')
                    print(f"Charged: {charged}")
                    print(f"Drained: {drained}")
                    print(f"Current level: Calculated from readings")
            else:
                print("No Body Battery data available")
        except Exception as e:
            print(f"Error fetching Body Battery: {e}")
            print("Note: Body Battery requires compatible Garmin device")

        # Fetch stress data
        print_section(f"Stress Data - {today}")
        try:
            stress = client.get_stress_data(today.isoformat())
            if stress:
                print(f"Overall stress level: {stress.get('overallStressLevel', 'N/A')}")
                print(f"Rest stress duration (sec): {stress.get('restStressDuration', 'N/A')}")
                print(f"Activity stress duration (sec): {stress.get('activityStressDuration', 'N/A')}")
                print(f"Low stress duration (sec): {stress.get('lowStressDuration', 'N/A')}")
                print(f"Medium stress duration (sec): {stress.get('mediumStressDuration', 'N/A')}")
                print(f"High stress duration (sec): {stress.get('highStressDuration', 'N/A')}")
        except Exception as e:
            print(f"Error fetching stress data: {e}")

        print_section("POC Complete!")
        print("✅ Successfully connected to Garmin and fetched your data!")
        print("\nNext steps:")
        print("1. Review the data above to confirm accuracy")
        print("2. If everything looks good, proceed with integration")
        print("3. Consider implementing the wearable-agnostic interface")

    except GarminConnectAuthenticationError:
        print("\n❌ Authentication failed!")
        print("Please check your email and password.")
        print("\nTroubleshooting:")
        print("- Make sure you can log into https://connect.garmin.com with these credentials")
        print("- Check if 2FA is enabled (may cause issues)")
        print("- Try again in a few minutes if you've had multiple failed attempts")
        sys.exit(1)

    except GarminConnectConnectionError as e:
        print(f"\n❌ Connection error: {e}")
        print("Please check your internet connection and try again.")
        sys.exit(1)

    except Exception as e:
        print(f"\n❌ Unexpected error: {e}")
        print(f"Error type: {type(e).__name__}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()

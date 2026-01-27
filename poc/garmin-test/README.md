# Garmin Connect POC

A simple proof-of-concept CLI tool to test the unofficial `python-garminconnect` library and validate that we can fetch your personal Garmin data.

## Purpose

Before integrating Garmin data fetching into the main health assistant system, this POC verifies:
- ✅ Authentication works with your Garmin credentials
- ✅ We can fetch sleep, HRV, activities, stress, and Body Battery data
- ✅ The data quality is sufficient for our causal modeling needs

## Prerequisites

- Python 3.11+
- [uv](https://github.com/astral-sh/uv) (modern Python package manager)
- Active Garmin Connect account
- Data synced to Garmin Connect (from your watch)

## Quick Start

### 1. Install uv (if not already installed)

```bash
# macOS/Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Or with pip
pip install uv
```

### 2. Run the POC

```bash
cd poc/garmin-test

# uv will automatically create a virtual environment and install dependencies
uv run garmin_poc.py
```

That's it! `uv` handles everything automatically.

You'll be prompted for your Garmin Connect credentials:
```
Email: your_email@example.com
Password: ******** (hidden)
```

## What It Fetches

The POC will attempt to fetch and display:

1. **User Profile** - Your Garmin name
2. **Daily Stats** (today) - Steps, calories, distance, HR
3. **Sleep Data** (yesterday) - Sleep stages, duration, sleep score
4. **Recent Activities** - Last 5 workouts with details
5. **HRV Data** (yesterday) - Heart rate variability (if available)
6. **Body Battery** (today) - Energy levels (if available)
7. **Stress Data** (today) - Stress levels throughout the day

## Expected Output

```
============================================================
  Garmin Connect POC - Data Fetcher
============================================================

Enter your Garmin Connect credentials:
Email: your_email@example.com
Password:

🔄 Connecting to Garmin Connect...
✅ Login successful!

============================================================
  User Profile
============================================================

Name: John Doe

============================================================
  Daily Stats - 2026-01-27
============================================================

Steps: 8542
Calories: 2341
Distance (meters): 6234
...

[more sections...]

============================================================
  POC Complete!
============================================================

✅ Successfully connected to Garmin and fetched your data!
```

## Troubleshooting

### Authentication Failed

**Symptoms**:
```
❌ Authentication failed!
Please check your email and password.
```

**Solutions**:
1. Verify credentials work at https://connect.garmin.com
2. Check if you have 2FA enabled (may cause issues with unofficial API)
3. Wait a few minutes if you've had multiple failed login attempts
4. Make sure you're using your email, not username

### No Data Available

**Symptoms**:
```
Error fetching sleep data: ...
No HRV data available
```

**Possible Causes**:
1. **Device not synced** - Sync your watch with Garmin Connect app
2. **Old device** - HRV and Body Battery require newer Garmin devices (Fenix 6+, Forerunner 945+, etc.)
3. **Time zone** - Sleep data is usually available for yesterday, not today
4. **Data not collected** - Make sure you wore your watch overnight

### Connection Errors

**Symptoms**:
```
❌ Connection error: ...
```

**Solutions**:
1. Check internet connection
2. Verify Garmin Connect is up (not under maintenance)
3. Try again in a few minutes

## Security Note

⚠️ **This uses the unofficial Garmin API which is against Garmin's Terms of Service.**

For this POC:
- ✅ Safe for personal testing
- ✅ Low risk of account issues with occasional use
- ⚠️ Don't run this continuously or at high frequency
- ⚠️ Don't share your credentials with anyone

For production:
- Consider switching to Oura Ring (official personal API)
- Or use Terra API (approved aggregator)
- See `docs/wearable-data-sources.md` for alternatives

## Data Quality Assessment

After running the POC, assess:

1. **Completeness**: Do you see all expected data types?
2. **Accuracy**: Compare values with Garmin Connect app
3. **Timeliness**: How recent is the data?
4. **Device Support**: Does your device support HRV, Body Battery, etc.?

## Alternative: Environment Variables

For repeated testing without entering credentials each time:

```bash
export GARMIN_EMAIL="your_email@example.com"
export GARMIN_PASSWORD="your_password"
uv run garmin_poc.py
```

(You'll need to modify the script to check for these environment variables)

## Advanced: Manual Environment Management

If you want more control:

```bash
# Create virtual environment
uv venv

# Activate it
source .venv/bin/activate  # On Windows: .venv\Scripts\activate

# Install dependencies
uv pip install -e .

# Run script
python garmin_poc.py
```

## Next Steps

If POC is successful:

1. ✅ Confirms unofficial API works with your account
2. → Design wearable-agnostic interface (see `docs/wearable-data-sources.md`)
3. → Build Python microservice for production use
4. → Integrate with Go ingestion-service
5. → Set up hourly polling (low frequency to avoid issues)
6. → Consider long-term alternatives (Oura, Terra)

## Resources

- [python-garminconnect GitHub](https://github.com/cyberjunky/python-garminconnect)
- [uv Documentation](https://docs.astral.sh/uv/)
- [Garmin Connect](https://connect.garmin.com)
- [Alternative: garth library](https://github.com/matin/garth)

---

**Last Updated**: January 2026
**Status**: POC Ready for Testing

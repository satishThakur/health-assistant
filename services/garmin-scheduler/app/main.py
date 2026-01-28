"""FastAPI application for Garmin scheduler service."""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.responses import JSONResponse

from .config import settings
from .scheduler import GarminSyncScheduler

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)

# Global scheduler instance
scheduler: GarminSyncScheduler = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan context manager for startup and shutdown."""
    global scheduler

    # Startup
    logger.info("Starting Garmin Scheduler Service")
    logger.info(f"Ingestion service URL: {settings.INGESTION_SERVICE_URL}")
    logger.info(f"Sync schedule: {settings.SYNC_CRON_HOUR}:{settings.SYNC_CRON_MINUTE}")

    scheduler = GarminSyncScheduler()
    scheduler.start()

    yield

    # Shutdown
    logger.info("Shutting down Garmin Scheduler Service")
    if scheduler:
        await scheduler.stop()


# Create FastAPI app
app = FastAPI(
    title="Garmin Scheduler Service",
    description="Scheduled service for syncing Garmin data to ingestion service",
    version="1.0.0",
    lifespan=lifespan,
)


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return JSONResponse(
        content={
            "status": "healthy",
            "service": "garmin-scheduler",
            "ingestion_url": settings.INGESTION_SERVICE_URL,
            "sync_schedule": f"{settings.SYNC_CRON_HOUR}:{settings.SYNC_CRON_MINUTE}",
        }
    )


@app.post("/sync/trigger")
async def trigger_sync():
    """Manually trigger a data sync (for testing)."""
    global scheduler

    if not scheduler:
        return JSONResponse(
            status_code=500,
            content={"status": "error", "message": "Scheduler not initialized"},
        )

    try:
        # Trigger sync in background
        await scheduler.trigger_manual_sync()
        return JSONResponse(
            content={"status": "success", "message": "Sync triggered successfully"}
        )
    except Exception as e:
        logger.error(f"Error triggering manual sync: {e}")
        return JSONResponse(
            status_code=500,
            content={"status": "error", "message": str(e)},
        )


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "service": "Garmin Scheduler",
        "status": "running",
        "endpoints": {
            "health": "/health",
            "trigger_sync": "/sync/trigger (POST)",
        },
    }

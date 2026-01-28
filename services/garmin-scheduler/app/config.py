"""Configuration management using Pydantic Settings."""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # Garmin credentials
    GARMIN_EMAIL: str
    GARMIN_PASSWORD: str

    # User configuration
    DEFAULT_USER_ID: str = "00000000-0000-0000-0000-000000000001"

    # Ingestion service
    INGESTION_SERVICE_URL: str = "http://ingestion-service:8083"

    # Scheduler configuration
    SYNC_CRON_HOUR: str = "*"  # Every hour by default
    SYNC_CRON_MINUTE: str = "0"  # At minute 0

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = True


# Global settings instance
settings = Settings()

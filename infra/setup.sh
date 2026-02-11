#!/bin/bash

# Setup script for Garmin data ingestion
# Usage: ./setup.sh

set -e

echo "🏥 Health Assistant - Setup"
echo "============================"
echo

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created"
    echo
    echo "⚠️  IMPORTANT: Please edit .env and add your Garmin credentials:"
    echo "   nano .env"
    echo
    echo "Update these values:"
    echo "   GARMIN_EMAIL=your_email@example.com"
    echo "   GARMIN_PASSWORD=your_password"
    echo
    read -p "Press Enter after updating .env file..."
else
    echo "✅ .env file already exists"
fi

# Check if credentials are set
if grep -q "your_email@example.com" .env; then
    echo "❌ Please update GARMIN_EMAIL in .env file"
    exit 1
fi

if grep -q "your_password" .env; then
    echo "❌ Please update GARMIN_PASSWORD in .env file"
    exit 1
fi

echo "✅ Credentials configured"
echo

# Start containers
echo "🚀 Starting Docker containers..."
docker compose down > /dev/null 2>&1 || true
docker compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 5

# Check health
echo "🏥 Checking service health..."
curl -s http://localhost:8083/health | jq '.' || echo "Ingestion service not ready yet"
curl -s http://localhost:8085/health | jq '.' || echo "Scheduler service not ready yet"

echo
echo "✅ Setup complete!"
echo
echo "📊 Next steps:"
echo "   1. Trigger initial sync:"
echo "      curl -X POST http://localhost:8085/sync/trigger"
echo
echo "   2. View dashboard data:"
echo "      curl http://localhost:8083/api/v1/dashboard/today | jq ."
echo
echo "   3. Run Flutter app:"
echo "      cd ../mobile_app && ./run.sh"
echo
echo "   4. View logs:"
echo "      docker compose logs -f garmin-scheduler"
echo

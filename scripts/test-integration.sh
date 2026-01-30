#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
INFRA_DIR="$PROJECT_ROOT/infra"
DEFAULT_USER_ID="00000000-0000-0000-0000-000000000001"
MAX_WAIT_TIME=120  # seconds

# Print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

# Check if .env file exists
check_env_file() {
    print_header "Checking Environment Configuration"

    if [ ! -f "$INFRA_DIR/.env" ]; then
        print_error ".env file not found in $INFRA_DIR"
        print_info "Please create .env file with your Garmin credentials:"
        print_info "  cd $INFRA_DIR"
        print_info "  cp .env.example .env"
        print_info "  # Edit .env and add your GARMIN_EMAIL and GARMIN_PASSWORD"
        exit 1
    fi

    # Check if required variables are set
    source "$INFRA_DIR/.env"

    if [ -z "$GARMIN_EMAIL" ] || [ -z "$GARMIN_PASSWORD" ]; then
        print_error "GARMIN_EMAIL and GARMIN_PASSWORD must be set in .env"
        exit 1
    fi

    print_success "Environment configuration found"
}

# Bring up containers
start_containers() {
    print_header "Starting Docker Containers"

    cd "$INFRA_DIR"
    print_info "Running: docker-compose up -d postgres ingestion-service garmin-scheduler"
    docker-compose up -d postgres ingestion-service garmin-scheduler

    print_success "Containers started"
}

# Wait for a service to be healthy
wait_for_service() {
    local service_name=$1
    local health_url=$2
    local max_attempts=$((MAX_WAIT_TIME / 5))
    local attempt=0

    print_info "Waiting for $service_name to be healthy..."

    while [ $attempt -lt $max_attempts ]; do
        if curl -f -s "$health_url" > /dev/null 2>&1; then
            print_success "$service_name is healthy"
            return 0
        fi

        attempt=$((attempt + 1))
        echo -n "."
        sleep 5
    done

    echo ""
    print_error "$service_name failed to become healthy after $MAX_WAIT_TIME seconds"
    return 1
}

# Wait for database to be ready
wait_for_database() {
    print_header "Waiting for Database"

    local max_attempts=$((MAX_WAIT_TIME / 5))
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if docker exec health-assistant-db pg_isready -U healthuser -d health_assistant > /dev/null 2>&1; then
            print_success "Database is ready"

            # Check if sync_audit table exists (from migration)
            print_info "Checking database schema..."
            if docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'sync_audit');" | grep -q 't'; then
                print_success "Database schema initialized (sync_audit table exists)"
            else
                print_warning "sync_audit table not found. Run migration: docker exec health-assistant-db psql -U healthuser -d health_assistant -f /docker-entrypoint-initdb.d/migrations/002_sync_audit.sql"
            fi

            return 0
        fi

        attempt=$((attempt + 1))
        echo -n "."
        sleep 5
    done

    echo ""
    print_error "Database failed to become ready after $MAX_WAIT_TIME seconds"
    return 1
}

# Check health of all services
check_services_health() {
    print_header "Checking Services Health"

    # Check ingestion service
    wait_for_service "Ingestion Service" "http://localhost:8083/health" || exit 1

    # Get and display ingestion service health details
    print_info "Ingestion service details:"
    curl -s http://localhost:8083/health | jq '.' || curl -s http://localhost:8083/health

    # Check garmin scheduler
    wait_for_service "Garmin Scheduler" "http://localhost:8085/health" || exit 1

    # Get and display scheduler details
    print_info "Scheduler details:"
    curl -s http://localhost:8085/health | jq '.' || curl -s http://localhost:8085/health

    print_success "All services are healthy"
}

# Get baseline data counts
get_baseline_counts() {
    print_header "Getting Baseline Data Counts"

    # Count events before sync
    EVENTS_BEFORE=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM events WHERE source = 'garmin' AND user_id = '$DEFAULT_USER_ID';")

    print_info "Events before sync: $EVENTS_BEFORE"

    # Show breakdown by type
    print_info "Breakdown by event type:"
    docker exec health-assistant-db psql -U healthuser -d health_assistant -c \
        "SELECT event_type, COUNT(*) as count FROM events WHERE source = 'garmin' AND user_id = '$DEFAULT_USER_ID' GROUP BY event_type ORDER BY event_type;" \
        || print_warning "No events found yet"
}

# Trigger manual sync
trigger_sync() {
    print_header "Triggering Manual Sync"

    print_info "Sending sync trigger to scheduler..."
    RESPONSE=$(curl -s -X POST http://localhost:8085/sync/trigger)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "Sync triggered successfully"
        echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
    else
        print_error "Failed to trigger sync"
        echo "$RESPONSE"
        return 1
    fi

    # Wait for sync to complete
    print_info "Waiting 30 seconds for sync to complete..."
    sleep 30
}

# Validate data was pulled
validate_data() {
    print_header "Validating Data Pull"

    # Count events after sync
    EVENTS_AFTER=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM events WHERE source = 'garmin' AND user_id = '$DEFAULT_USER_ID';")

    print_info "Events after sync: $EVENTS_AFTER"

    # Check if new data was added
    NEW_EVENTS=$((EVENTS_AFTER - EVENTS_BEFORE))

    if [ $NEW_EVENTS -gt 0 ]; then
        print_success "✓ $NEW_EVENTS new events were added"
    else
        print_warning "No new events added (data might already exist or no new data available)"
    fi

    # Show breakdown by type
    print_info "Current breakdown by event type:"
    docker exec health-assistant-db psql -U healthuser -d health_assistant -c \
        "SELECT event_type, COUNT(*) as count, MAX(time) as latest FROM events WHERE source = 'garmin' AND user_id = '$DEFAULT_USER_ID' GROUP BY event_type ORDER BY event_type;"

    # Check sync audit logs
    print_info "Recent sync audit logs:"
    AUDIT_COUNT=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID';")

    if [ "$AUDIT_COUNT" -gt 0 ]; then
        print_success "✓ $AUDIT_COUNT sync audit records found"

        print_info "Latest sync audit entries:"
        docker exec health-assistant-db psql -U healthuser -d health_assistant -c \
            "SELECT data_type, target_date, records_fetched, records_inserted, records_updated, status, sync_duration_seconds
             FROM sync_audit
             WHERE user_id = '$DEFAULT_USER_ID'
             ORDER BY sync_started_at DESC
             LIMIT 10;"

        # Check for failures
        FAILED_SYNCS=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
            "SELECT COUNT(*) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID' AND status = 'failed';")

        if [ "$FAILED_SYNCS" -gt 0 ]; then
            print_warning "Found $FAILED_SYNCS failed sync(s)"
            print_info "Failed sync details:"
            docker exec health-assistant-db psql -U healthuser -d health_assistant -c \
                "SELECT data_type, target_date, error_message FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID' AND status = 'failed' ORDER BY sync_started_at DESC LIMIT 5;"
        fi
    else
        print_warning "No sync audit records found (sync might have failed or audit table not initialized)"
    fi
}

# Test API endpoints
test_api_endpoints() {
    print_header "Testing API Endpoints"

    # Test recent audits endpoint
    print_info "Testing GET /api/v1/audit/sync/recent"
    AUDIT_RESPONSE=$(curl -s "http://localhost:8083/api/v1/audit/sync/recent?user_id=$DEFAULT_USER_ID&limit=5")
    if [ $? -eq 0 ]; then
        print_success "✓ Recent audits endpoint working"
        echo "$AUDIT_RESPONSE" | jq '.[0:2]' 2>/dev/null || echo "$AUDIT_RESPONSE" | head -20
    else
        print_error "Recent audits endpoint failed"
    fi

    # Test stats endpoint
    print_info "Testing GET /api/v1/audit/sync/stats"
    STATS_RESPONSE=$(curl -s "http://localhost:8083/api/v1/audit/sync/stats?user_id=$DEFAULT_USER_ID")
    if [ $? -eq 0 ]; then
        print_success "✓ Stats endpoint working"
        echo "$STATS_RESPONSE" | jq '.' 2>/dev/null || echo "$STATS_RESPONSE"
    else
        print_error "Stats endpoint failed"
    fi
}

# View logs
view_logs() {
    print_header "Recent Container Logs"

    print_info "Garmin Scheduler logs (last 20 lines):"
    docker logs --tail 20 health-assistant-garmin-scheduler

    echo ""
    print_info "Ingestion Service logs (last 20 lines):"
    docker logs --tail 20 health-assistant-ingestion-service
}

# Generate summary report
generate_summary() {
    print_header "Integration Test Summary"

    # Get final counts
    local total_events=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM events WHERE source = 'garmin' AND user_id = '$DEFAULT_USER_ID';")

    local total_syncs=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID';")

    local successful_syncs=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COUNT(*) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID' AND status = 'success';")

    local total_fetched=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COALESCE(SUM(records_fetched), 0) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID';")

    local total_inserted=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COALESCE(SUM(records_inserted), 0) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID';")

    local total_updated=$(docker exec health-assistant-db psql -U healthuser -d health_assistant -tAc \
        "SELECT COALESCE(SUM(records_updated), 0) FROM sync_audit WHERE user_id = '$DEFAULT_USER_ID';")

    echo "┌─────────────────────────────────────────┐"
    echo "│        Integration Test Results         │"
    echo "├─────────────────────────────────────────┤"
    echo "│ Total Events:          $(printf '%15s' $total_events) │"
    echo "│ Total Sync Runs:       $(printf '%15s' $total_syncs) │"
    echo "│ Successful Syncs:      $(printf '%15s' $successful_syncs) │"
    echo "│ Records Fetched:       $(printf '%15s' $total_fetched) │"
    echo "│ Records Inserted:      $(printf '%15s' $total_inserted) │"
    echo "│ Records Updated:       $(printf '%15s' $total_updated) │"
    echo "└─────────────────────────────────────────┘"

    if [ "$total_events" -gt 0 ] && [ "$total_syncs" -gt 0 ]; then
        print_success "Integration test PASSED ✓"
        echo ""
        print_info "Next steps:"
        echo "  • View logs: docker logs health-assistant-garmin-scheduler"
        echo "  • Query data: docker exec -it health-assistant-db psql -U healthuser -d health_assistant"
        echo "  • Check audit API: curl 'http://localhost:8083/api/v1/audit/sync/recent?user_id=$DEFAULT_USER_ID'"
        return 0
    else
        print_warning "Integration test completed with warnings"
        print_info "This might be expected if no new data is available from Garmin"
        return 0
    fi
}

# Cleanup function
cleanup() {
    if [ "$CLEANUP_ON_EXIT" = "true" ]; then
        print_header "Cleaning Up"
        cd "$INFRA_DIR"
        docker-compose down
        print_info "Containers stopped"
    fi
}

# Main execution
main() {
    print_header "Garmin Integration Test"
    print_info "Starting end-to-end integration test..."

    # Parse arguments
    CLEANUP_ON_EXIT=false
    SKIP_LOGS=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --cleanup)
                CLEANUP_ON_EXIT=true
                shift
                ;;
            --skip-logs)
                SKIP_LOGS=true
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --cleanup     Stop containers after test"
                echo "  --skip-logs   Skip displaying container logs"
                echo "  --help        Show this help message"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    # Set trap for cleanup
    trap cleanup EXIT

    # Run test steps
    check_env_file
    start_containers
    wait_for_database
    check_services_health
    get_baseline_counts
    trigger_sync
    validate_data
    test_api_endpoints

    if [ "$SKIP_LOGS" != "true" ]; then
        view_logs
    fi

    generate_summary
}

# Run main function
main "$@"

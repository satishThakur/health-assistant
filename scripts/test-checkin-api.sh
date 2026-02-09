#!/bin/bash

# Test script for Daily Check-in API endpoints
# Usage: ./test-checkin-api.sh [base_url]
# Example: ./test-checkin-api.sh http://localhost:8083

set -e

BASE_URL="${1:-http://localhost:8083}"

echo "========================================="
echo "Testing Daily Check-in API"
echo "Base URL: $BASE_URL"
echo "========================================="
echo

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Submit check-in
echo -e "${BLUE}1. Testing POST /api/v1/checkin${NC}"
echo "Submitting daily check-in..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/checkin" \
  -H "Content-Type: application/json" \
  -d '{
    "energy": 8,
    "mood": 7,
    "focus": 9,
    "physical": 7,
    "notes": "Felt great after morning run"
  }')

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

# Test 2: Get latest check-in
echo -e "${BLUE}2. Testing GET /api/v1/checkin/latest${NC}"
echo "Fetching today's check-in..."
RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/checkin/latest")

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

# Test 3: Get check-in history
echo -e "${BLUE}3. Testing GET /api/v1/checkin/history?days=7${NC}"
echo "Fetching 7-day check-in history..."
RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/checkin/history?days=7")

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

# Test 4: Get today's dashboard
echo -e "${BLUE}4. Testing GET /api/v1/dashboard/today${NC}"
echo "Fetching today's dashboard (check-in + Garmin data)..."
RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/dashboard/today")

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

# Test 5: Get week trends
echo -e "${BLUE}5. Testing GET /api/v1/trends/week${NC}"
echo "Fetching 7-day trends..."
RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/trends/week")

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

# Test 6: Get correlations
echo -e "${BLUE}6. Testing GET /api/v1/insights/correlations?days=30${NC}"
echo "Calculating correlations (requires 30+ days of data)..."
RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/insights/correlations?days=30")

echo "Response:"
echo "$RESPONSE" | jq '.'
echo

echo "========================================="
echo -e "${GREEN}✓ All tests completed!${NC}"
echo "========================================="
echo

# Test invalid payloads
echo -e "${YELLOW}Testing validation (invalid payloads):${NC}"
echo

echo -e "${BLUE}7. Testing validation - energy out of range${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/checkin" \
  -H "Content-Type: application/json" \
  -d '{
    "energy": 11,
    "mood": 7,
    "focus": 8,
    "physical": 7
  }')

echo "Response (should fail):"
echo "$RESPONSE" | jq '.'
echo

echo -e "${BLUE}8. Testing validation - missing required field${NC}"
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/checkin" \
  -H "Content-Type: application/json" \
  -d '{
    "energy": 8,
    "mood": 7
  }')

echo "Response (should fail):"
echo "$RESPONSE" | jq '.'
echo

echo "========================================="
echo -e "${GREEN}✓ Validation tests completed!${NC}"
echo "========================================="

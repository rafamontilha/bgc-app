#!/bin/bash

# ============================================
# BGC - Testes E2E Docker Environment
# ============================================

echo "╔════════════════════════════════════════════════════════════╗"
echo "║  🧪 BGC - Testes E2E - Ambiente Docker                    ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

test_service() {
    local name=$1
    local url=$2
    local expected=$3

    echo -n "Testing $name... "

    response=$(curl -s -w "\n%{http_code}" -m 5 "$url" 2>&1)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        if [ -z "$expected" ] || echo "$body" | grep -q "$expected"; then
            echo -e "${GREEN}✓ PASSED${NC}"
            ((PASSED++))
            return 0
        else
            echo -e "${RED}✗ FAILED${NC} (content mismatch)"
            ((FAILED++))
            return 1
        fi
    else
        echo -e "${RED}✗ FAILED${NC} (HTTP $http_code)"
        ((FAILED++))
        return 1
    fi
}

test_api_simulation() {
    local ncm=$1
    local description=$2

    echo ""
    echo -e "${BLUE}Testing API Simulation: $description${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    payload="{\"ncm\":\"$ncm\",\"volume_kg\":10000,\"max_results\":3}"

    echo "Request: POST /v1/simulator/destinations"
    echo "Payload: $payload"
    echo ""

    response=$(curl -s -w "\n%{http_code}" -m 10 \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$payload" \
        "http://localhost:8080/v1/simulator/destinations" 2>&1)

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        echo -e "${GREEN}✓ API Response: SUCCESS${NC}"
        echo ""

        # Parse JSON response
        destinations=$(echo "$body" | grep -o '"destinations":\[' | wc -l)
        ncm_result=$(echo "$body" | grep -o "\"ncm\":\"[^\"]*\"" | head -1 | cut -d'"' -f4)
        processing_time=$(echo "$body" | grep -o '"processing_time_ms":[0-9]*' | cut -d':' -f2)

        echo "📊 Results:"
        echo "  • NCM: $ncm_result"
        echo "  • Processing Time: ${processing_time}ms"
        echo "  • Response Size: $(echo "$body" | wc -c) bytes"
        echo ""

        # Extract first destination
        first_country=$(echo "$body" | grep -o '"country_name":"[^"]*"' | head -1 | cut -d'"' -f4)
        first_score=$(echo "$body" | grep -o '"score":[0-9.]*' | head -1 | cut -d':' -f2)

        if [ ! -z "$first_country" ]; then
            echo "🌍 Top Destination:"
            echo "  • Country: $first_country"
            echo "  • Score: $first_score"
        fi

        ((PASSED++))
        return 0
    elif [ "$http_code" = "429" ]; then
        echo -e "${YELLOW}⚠ RATE LIMIT EXCEEDED${NC}"
        echo "Message: $(echo "$body" | grep -o '"message":"[^"]*"' | cut -d'"' -f4)"
        ((FAILED++))
        return 1
    elif [ "$http_code" = "404" ]; then
        echo -e "${RED}✗ NCM NOT FOUND${NC}"
        echo "Message: $(echo "$body" | grep -o '"message":"[^"]*"' | cut -d'"' -f4)"
        ((FAILED++))
        return 1
    else
        echo -e "${RED}✗ API ERROR (HTTP $http_code)${NC}"
        echo "Response: $body"
        ((FAILED++))
        return 1
    fi
}

# ============================================
# Test Suite
# ============================================

echo "📋 Test Suite: Infrastructure"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
test_service "API Health" "http://localhost:8080/healthz" "status"
test_service "Frontend Home" "http://localhost:3000/" "BGC"
test_service "Simulator Page" "http://localhost:3000/simulator" "Simulador"

echo ""
echo "📋 Test Suite: API Integration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test with valid NCM (Sugar)
test_api_simulation "17011400" "Açúcar de cana (Valid)"

# Test with valid NCM (Soy)
test_api_simulation "12010090" "Sementes de soja (Valid)"

# Test with valid NCM (Iron Ore)
test_api_simulation "26011200" "Minério de ferro (Valid)"

# Test with invalid NCM
echo ""
echo -e "${BLUE}Testing API: Invalid NCM${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
response=$(curl -s -w "\n%{http_code}" -m 5 \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"ncm":"99999999","volume_kg":10000}' \
    "http://localhost:8080/v1/simulator/destinations" 2>&1)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "404" ] || [ "$http_code" = "400" ]; then
    echo -e "${GREEN}✓ Correctly rejected invalid NCM${NC}"
    echo "Error message: $(echo "$body" | grep -o '"message":"[^"]*"' | cut -d'"' -f4)"
    ((PASSED++))
else
    echo -e "${RED}✗ Should reject invalid NCM${NC}"
    ((FAILED++))
fi

# ============================================
# Summary
# ============================================

echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║  📊 Test Results Summary                                  ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""
echo -e "  ${GREEN}Passed: $PASSED${NC}"
echo -e "  ${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ Some tests failed. Check logs above.${NC}"
    exit 1
fi

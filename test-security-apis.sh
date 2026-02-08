#!/bin/bash
# Test script for x402 Security APIs
# Usage: ./test-security-apis.sh [base-url]

set -e

BASE_URL="${1:-http://localhost:8080}"

echo "🧪 Testing x402 Security APIs"
echo "=============================="
echo "Base URL: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper function to make requests
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expect_status="$5"
    
    echo "Testing: $name"
    echo "  $method $endpoint"
    
    if [ "$method" = "POST" ] && [ -n "$data" ]; then
        HTTP_STATUS=$(curl -s -o /tmp/response.json -w "%{http_code}" \
            -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint" 2>/dev/null || echo "000")
    else
        HTTP_STATUS=$(curl -s -o /tmp/response.json -w "%{http_code}" \
            -X "$method" \
            "$BASE_URL$endpoint" 2>/dev/null || echo "000")
    fi
    
    if [ "$HTTP_STATUS" = "$expect_status" ]; then
        echo -e "  ${GREEN}✓${NC} Got expected status $expect_status"
        ((TESTS_PASSED++))
    else
        echo -e "  ${RED}✗${NC} Expected $expect_status, got $HTTP_STATUS"
        ((TESTS_FAILED++))
    fi
    
    if [ -f /tmp/response.json ] && [ -s /tmp/response.json ]; then
        echo "  Response preview:"
        cat /tmp/response.json | head -c 200
        echo "..."
    fi
    echo ""
}

# Test 1: Health check
echo "1. Testing Health Endpoint (Free)"
test_endpoint "Health" "GET" "/health" "" "200"

# Test 2: Service info
echo "2. Testing Service Info (Free)"
test_endpoint "Service Info" "GET" "/" "" "200"

# Test 3: x402 config
echo "3. Testing x402 Config (Free)"
test_endpoint "x402 Config" "GET" "/.well-known/x402" "" "200"

# Test 4: Contract scanner without payment (should return 402)
echo "4. Testing Contract Scanner (Payment Required)"
test_endpoint "Contract Scan (no payment)" "POST" "/api/scan-contract" '{"address":"0x4200000000000000000000000000000000000006","chain":"base"}' "402"

# Test 5: Agent score without payment
echo "5. Testing Agent Score (Payment Required)"
test_endpoint "Agent Score (no payment)" "POST" "/api/agent-score" '{"agent_id":"1941"}' "402"

# Test 6: TX preflight without payment
echo "6. Testing TX Pre-flight (Payment Required)"
test_endpoint "TX Pre-flight (no payment)" "POST" "/api/tx-preflight" '{"from":"0x...","to":"0x...","value":"0x0","data":"0x"}' "402"

# Test 7: Prompt test without payment
echo "7. Testing Prompt Injection (Payment Required)"
test_endpoint "Prompt Test (no payment)" "POST" "/api/prompt-test" '{"prompt":"Hello world"}' "402"

# Test 8: Gas price without payment
echo "8. Testing Gas Price (Payment Required)"
test_endpoint "Gas Price (no payment)" "GET" "/api/gas" "" "402"

# Test 9: Test invalid method
echo "9. Testing Invalid Method"
test_endpoint "Contract Scan (GET method)" "GET" "/api/scan-contract" "" "405"

# Test 10: Test missing fields
echo "10. Testing Missing Fields"
test_endpoint "Contract Scan (missing address)" "POST" "/api/scan-contract" '{"chain":"base"}' "400"

echo "=============================="
echo "Test Results:"
echo -e "  ${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "  ${RED}Failed: $TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi

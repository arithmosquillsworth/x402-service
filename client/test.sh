#!/bin/bash
# Test script for x402 service endpoints

set -e

API_URL="${X402_API_URL:-http://localhost:8080}"

echo "🔍 Testing x402 Service at $API_URL"
echo "========================================="
echo

# Test health endpoint
echo "1. Testing /health (free)"
curl -s "$API_URL/health" | jq .
echo

# Test config endpoint
echo "2. Testing /.well-known/x402 (free)"
curl -s "$API_URL/.well-known/x402" | jq .
echo

# Test gas endpoint without payment (expect 402)
echo "3. Testing /api/gas without payment (expect 402)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$API_URL/api/gas" | jq .
echo

# Test validators endpoint without payment (expect 402)
echo "4. Testing /api/validators without payment (expect 402)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$API_URL/api/validators" | jq .
echo

# Test price endpoint without payment (expect 402)
echo "5. Testing /api/price without payment (expect 402)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$API_URL/api/price" | jq .
echo

echo "========================================="
echo "✅ Tests complete"
echo
echo "To test paid endpoints, you need an x402 payment token."
echo "See: https://github.com/arithmosquillsworth/x402-service#usage"

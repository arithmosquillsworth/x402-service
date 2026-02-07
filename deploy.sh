#!/bin/bash
# Quick deployment script for x402-service
# This script deploys the x402 service to a publicly accessible URL

set -e

echo "🚀 x402 Service Deployment Script"
echo "=================================="

# Check if binary exists
if [ ! -f "./x402-service" ]; then
    echo "❌ x402-service binary not found. Building..."
    go build -o x402-service main.go
fi

# Option 1: Try to use cloudflared if configured
if command -v cloudflared &> /dev/null; then
    echo "📡 Attempting Cloudflare tunnel..."
    cloudflared tunnel --url http://localhost:8080 &
    CLOUDFLARE_PID=$!
fi

# Option 2: Try to use SSH reverse tunnel if VPS configured
# ssh -R 8080:localhost:8080 user@vps.example.com &

# Start the service
echo "🎯 Starting x402 service on port 8080..."
export RECEIVER_ADDRESS="0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91"
export PORT="8080"
export ETH_RPC_URL="https://eth.drpc.org"

./x402-service &
SERVICE_PID=$!

echo "✅ Service started with PID: $SERVICE_PID"
echo ""
echo "Service endpoints:"
echo "  Health:    http://localhost:8080/health"
echo "  Config:    http://localhost:8080/.well-known/x402"
echo "  Gas:       http://localhost:8080/api/gas (0.001 USDC)"
echo "  Price:     http://localhost:8080/api/price (0.002 USDC)"
echo "  Validators: http://localhost:8080/api/validators (0.005 USDC)"
echo ""
echo "Press Ctrl+C to stop"

# Wait for interrupt
trap "echo 'Stopping...'; kill $SERVICE_PID $CLOUDFLARE_PID 2>/dev/null; exit" INT
wait

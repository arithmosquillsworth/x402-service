#!/bin/bash
# Deploy x402 Security API Service to Render
# Usage: ./deploy-render.sh [service-name]

set -e

SERVICE_NAME="${1:-x402-security-service}"

echo "🚀 Deploying x402 Security API Service to Render..."
echo ""

# Check if render CLI is installed
if ! command -v render &> /dev/null; then
    echo "❌ Render CLI not found. Install with:"
    echo "   curl -fsSL https://raw.githubusercontent.com/render-oss/render-cli/main/install.sh | bash"
    exit 1
fi

# Check if logged in
if ! render whoami &> /dev/null; then
    echo "🔑 Please login to Render:"
    render login
fi

# Check environment variables
echo "📋 Checking environment variables..."

if [ -z "$BASESCAN_API_KEY" ]; then
    echo "⚠️  Warning: BASESCAN_API_KEY not set"
    echo "   Contract scanning on Base will not work"
fi

if [ -z "$ETHERSCAN_API_KEY" ]; then
    echo "⚠️  Warning: ETHERSCAN_API_KEY not set"
    echo "   Contract scanning on Ethereum will not work"
fi

# Build locally to verify
echo ""
echo "🔨 Building locally to verify..."
go build -o x402-service .

if [ $? -ne 0 ]; then
    echo "❌ Build failed! Fix errors before deploying."
    exit 1
fi

echo "✅ Build successful"

# Deploy
echo ""
echo "📤 Deploying to Render..."
render deploy --service "$SERVICE_NAME"

echo ""
echo "✅ Deployment initiated!"
echo ""
echo "🔗 Useful commands:"
echo "   render services                    # List services"
echo "   render logs --service $SERVICE_NAME    # View logs"
echo "   render ssh --service $SERVICE_NAME     # SSH into service"
echo ""
echo "📝 Don't forget to set environment variables in Render dashboard:"
echo "   - BASESCAN_API_KEY"
echo "   - ETHERSCAN_API_KEY"

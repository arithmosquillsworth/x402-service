#!/bin/bash
# Render Deployment Script for x402-service
# Usage: ./deploy-render.sh [service-name]

set -e

SERVICE_NAME="${1:-x402-service}"
REPO_DIR="$HOME/projects/arithmosquillsworth/x402-service"

echo "🚀 Render Deployment Script"
echo "============================"
echo "Service: $SERVICE_NAME"
echo ""

# Check if we're in the right directory
if [ ! -f "$REPO_DIR/render.yaml" ]; then
    echo "❌ Error: render.yaml not found in $REPO_DIR"
    echo "Please run this script from the x402-service directory"
    exit 1
fi

cd "$REPO_DIR"

# Check if render CLI is available
if ! command -v render &> /dev/null; then
    echo "⚠️  Render CLI not found"
    echo "Install with: curl -fsSL https://raw.githubusercontent.com/render-oss/render-cli/main/install.sh | bash"
    echo ""
    echo "Alternatively, deploy manually:"
    echo "1. Go to https://dashboard.render.com"
    echo "2. Click 'New +' → 'Web Service'"
    echo "3. Connect GitHub repo: arithmosquillsworth/x402-service"
    echo "4. Select 'Docker' runtime"
    echo "5. Set environment variables if needed"
    exit 1
fi

# Check if logged in
echo "🔍 Checking Render login status..."
if ! render whoami &> /dev/null; then
    echo "❌ Not logged in to Render"
    echo "Run: render login"
    exit 1
fi

echo "✅ Logged in as: $(render whoami)"
echo ""

# Check if service already exists
echo "🔍 Checking if service exists..."
if render services list | grep -q "$SERVICE_NAME"; then
    echo "✅ Service exists, deploying latest commit..."
    render deploy "$SERVICE_NAME" --wait
else
    echo "🆕 Creating new service from render.yaml..."
    render blueprint apply render.yaml --wait
fi

echo ""
echo "✅ Deployment complete!"
echo ""
echo "📊 Service URL: https://$SERVICE_NAME.onrender.com"
echo "🔍 Check status: render services info $SERVICE_NAME"
echo "📜 View logs: render logs $SERVICE_NAME"

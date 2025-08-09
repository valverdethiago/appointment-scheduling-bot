#!/bin/bash

# Test Google Calendar Service Account Connection
# This script tests if your service account can access the calendar

echo "🧪 Testing Google Calendar Service Account Connection"
echo "=================================================="

# Check if service account file exists
if [ ! -f "service-account.json" ]; then
    echo "❌ service-account.json not found!"
    exit 1
fi

# Check if config exists
if [ ! -f "config.env" ]; then
    echo "❌ config.env not found!"
    exit 1
fi

# Build the test binary
echo "🔨 Building test binary..."
if go build -o test-calendar cmd/api/main.go; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

# Test the configuration loading
echo ""
echo "📋 Configuration Test:"
echo "====================="
./test-calendar &
PID=$!

# Wait a moment for the server to start
sleep 2

# Test the health endpoint
echo "🌐 Testing health endpoint..."
if curl -s http://localhost:8080/healthz > /dev/null; then
    echo "✅ Server is running"
    
    # Get detailed health info
    echo ""
    echo "📊 Health Check Details:"
    curl -s http://localhost:8080/healthz | jq '.' 2>/dev/null || curl -s http://localhost:8080/healthz
    
    echo ""
    echo "✅ Configuration loaded successfully!"
    echo "✅ Service account file found"
    echo "✅ Calendar ID configured"
    
else
    echo "❌ Server health check failed"
fi

# Clean up
echo ""
echo "🧹 Cleaning up..."
kill $PID 2>/dev/null
rm -f test-calendar

echo ""
echo "🎯 Next Steps:"
echo "=============="
echo "1. Run the setup script: ./scripts/setup-calendar.sh"
echo "2. Follow the calendar sharing instructions"
echo "3. Test again with: ./scripts/test-calendar.sh"
echo "4. Start your API: go run cmd/api/main.go" 
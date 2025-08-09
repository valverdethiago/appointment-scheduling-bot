#!/bin/bash

# Google Calendar Setup Script for Service Account
# This script helps you configure the necessary permissions for your service account

echo "🔧 Google Calendar Service Account Setup"
echo "========================================"

# Read the service account email from the JSON file
if [ -f "service-account.json" ]; then
    SERVICE_ACCOUNT_EMAIL=$(grep -o '"client_email": "[^"]*"' service-account.json | cut -d'"' -f4)
    echo "✅ Found service account: $SERVICE_ACCOUNT_EMAIL"
else
    echo "❌ service-account.json not found!"
    exit 1
fi

# Read the calendar ID from config
if [ -f "config.env" ]; then
    CALENDAR_ID=$(grep "GCAL_CALENDAR_ID=" config.env | cut -d'=' -f2)
    echo "✅ Calendar ID: $CALENDAR_ID"
else
    echo "❌ config.env not found!"
    exit 1
fi

echo ""
echo "📋 Next Steps for Test Account Setup:"
echo "====================================="
echo ""
echo "1. Go to Google Calendar: https://calendar.google.com"
echo "2. Sign in with: appointmentschedulingbot@gmail.com"
echo "3. In the left sidebar, find your calendar '$CALENDAR_ID'"
echo "4. Click the three dots next to the calendar name"
echo "5. Select 'Settings and sharing'"
echo "6. Scroll down to 'Share with specific people'"
echo "7. Click '+ Add people'"
echo "8. Add this service account email: $SERVICE_ACCOUNT_EMAIL"
echo "9. Set permission to 'Make changes to events'"
echo "10. Click 'Send'"
echo ""
echo "🔐 Alternative: Create a dedicated calendar"
echo "=========================================="
echo "1. In Google Calendar, click '+ Add other calendars'"
echo "2. Select 'Create new calendar'"
echo "3. Name it 'Appointment Scheduling Bot'"
echo "4. Share it with: $SERVICE_ACCOUNT_EMAIL"
echo "5. Update your config.env with the new calendar ID"
echo ""
echo "🧪 Test the connection:"
echo "======================"
echo "Run: go run scripts/test-calendar-connection.go"
echo "Or start the API: go run cmd/api/main.go"
echo "Then check: http://localhost:8080/healthz"
echo ""
echo "📚 Documentation:"
echo "================"
echo "Service Account Setup: https://developers.google.com/calendar/api/guides/auth"
echo "Calendar API: https://developers.google.com/calendar/api/v3/reference" 
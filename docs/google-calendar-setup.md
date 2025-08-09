# Google Calendar Service Account Setup

This document provides a complete guide for setting up Google Cloud service account authentication and Google Calendar integration for the Appointment Scheduling Bot.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Step 1: Create Google Cloud Project](#step-1-create-google-cloud-project)
- [Step 2: Enable Google Calendar API](#step-2-enable-google-calendar-api)
- [Step 3: Create Service Account](#step-3-create-service-account)
- [Step 4: Download Service Account Key](#step-4-download-service-account-key)
- [Step 5: Configure Calendar Permissions](#step-5-configure-calendar-permissions)
- [Step 6: Update Configuration](#step-6-update-configuration)
- [Step 7: Test the Setup](#step-7-test-the-setup)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)

## Overview

The Appointment Scheduling Bot uses Google Calendar API to manage appointments and check availability. To access the API, we use a service account with appropriate permissions. This setup allows the bot to:

- Read calendar availability
- Create new appointments
- Update existing appointments
- Delete appointments
- Manage calendar events programmatically

## Prerequisites

- Google Cloud account
- Access to Google Calendar
- Go development environment
- Basic understanding of Google Cloud Console

## Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click "Select a project" at the top of the page
3. Click "New Project"
4. Enter a project name (e.g., "Appointment Scheduling Bot")
5. Click "Create"

## Step 2: Enable Google Calendar API

1. In your new project, go to the [APIs & Services > Library](https://console.cloud.google.com/apis/library)
2. Search for "Google Calendar API"
3. Click on "Google Calendar API"
4. Click "Enable"

## Step 3: Create Service Account

1. Go to [APIs & Services > Credentials](https://console.cloud.google.com/apis/credentials)
2. Click "Create Credentials" > "Service Account"
3. Fill in the service account details:
   - **Name**: `appointment-scheduling-bot`
   - **Description**: `Service account for appointment scheduling bot`
4. Click "Create and Continue"
5. For roles, add:
   - **Basic** > **Editor** (for general project access)
6. Click "Continue"
7. Click "Done"

## Step 4: Download Service Account Key

1. In the service accounts list, click on your newly created service account
2. Go to the "Keys" tab
3. Click "Add Key" > "Create new key"
4. Choose "JSON" format
5. Click "Create"
6. The JSON file will automatically download
7. **Important**: Move this file to your project root and rename it to `service-account.json`
8. **Never commit this file to version control!**

## Step 5: Configure Calendar Permissions

### Option A: Use Existing Calendar (Recommended for Testing)

1. Go to [Google Calendar](https://calendar.google.com)
2. Sign in with the account that owns the calendar you want to use
3. In the left sidebar, find your calendar
4. Click the three dots next to the calendar name
5. Select "Settings and sharing"
6. Scroll down to "Share with specific people"
7. Click "+ Add people"
8. Add your service account email: `appointment-scheduling-bot@your-project-id.iam.gserviceaccount.com`
9. Set permission to **"Make changes to events"**
10. Click "Send"

### Option B: Create Dedicated Calendar (Recommended for Production)

1. In Google Calendar, click "+ Add other calendars"
2. Select "Create new calendar"
3. Name it "Appointment Scheduling Bot"
4. Set the timezone and other preferences
5. Click "Create calendar"
6. Follow steps 4-10 from Option A to share with the service account

## Step 6: Update Configuration

1. Open your `config.env` file
2. Update the following variables:

```env
# Google Calendar Configuration
GCAL_CALENDAR_ID=your-calendar-id@gmail.com
GOOGLE_CREDS_JSON=service-account.json
```

**Note**: Replace `your-calendar-id@gmail.com` with the actual calendar ID you're using.

## Step 7: Test the Setup

### Quick Test

Run the basic configuration test:

```bash
./scripts/test-calendar.sh
```

This will verify that:
- Configuration is loaded correctly
- Service account file is found
- Basic server functionality works

### Full API Test

Test the actual Google Calendar API connection:

```bash
go run scripts/test-calendar-connection.go
```

This will verify:
- Service account authentication
- Calendar access permissions
- Read/write capabilities

### Expected Output

A successful test should show:

```
🧪 Testing Google Calendar API Connection
========================================
✅ Configuration loaded
   Calendar ID: your-calendar-id@gmail.com
   Service Account: appointment-scheduling-bot@your-project-id.iam.gserviceaccount.com

🔌 Creating Google Calendar client...
✅ Client created successfully

📅 Testing calendar access...
✅ Successfully accessed calendar!
   Found X busy time blocks

📝 Testing event creation (optional)...
✅ Test event created successfully!
   Event ID: abc123...
🧹 Cleaning up test event...
✅ Test event cleaned up

🎉 Google Calendar Service Account Test Complete!
===============================================
✅ Configuration: Working
✅ Authentication: Working
✅ Calendar Access: Working
✅ API Connection: Working
✅ Write Permissions: Working
```

## Troubleshooting

### Common Issues

#### 1. "Failed to load config: GCAL_CALENDAR_ID is required"

**Solution**: Check your `config.env` file and ensure `GCAL_CALENDAR_ID` is set correctly.

#### 2. "Failed to get Google credentials"

**Solution**: Verify that `service-account.json` exists in your project root and has the correct format.

#### 3. "googleapi: Error 403: Forbidden"

**Solution**: The service account doesn't have permission to access the calendar. Follow Step 5 to share the calendar.

#### 4. "googleapi: Error 404: Not Found"

**Solution**: The calendar ID is incorrect or the service account doesn't have access. Verify the calendar ID and sharing permissions.

#### 5. "Failed to create JWT config"

**Solution**: The service account JSON file is corrupted or invalid. Re-download it from Google Cloud Console.

### Debug Steps

1. **Check Configuration**:
   ```bash
   cat config.env | grep GCAL
   ```

2. **Verify Service Account File**:
   ```bash
   ls -la service-account.json
   ```

3. **Test Basic Connection**:
   ```bash
   ./scripts/test-calendar.sh
   ```

4. **Check Calendar Sharing**:
   - Go to Google Calendar
   - Verify the service account email is listed in sharing settings
   - Ensure permissions are set to "Make changes to events"

## Security Best Practices

### 1. Service Account Key Security

- **Never commit** `service-account.json` to version control
- Store the file securely on your deployment server
- Use environment variables in production
- Rotate keys regularly

### 2. Calendar Permissions

- Grant minimal necessary permissions
- Use dedicated calendars for production
- Regularly review sharing settings
- Monitor API usage

### 3. Environment Configuration

- Use different service accounts for development and production
- Store sensitive configuration in environment variables
- Use `.env` files for local development only

### 4. API Rate Limits

- Be aware of Google Calendar API quotas
- Implement rate limiting in your application
- Monitor API usage in Google Cloud Console

## File Structure

After setup, your project should have:

```
appointment-scheduling-bot/
├── service-account.json          # Service account credentials (DO NOT COMMIT)
├── config.env                    # Configuration file
├── scripts/
│   ├── setup-calendar.sh        # Setup helper script
│   └── test-calendar.sh         # Basic test script
├── docs/
│   └── google-calendar-setup.md # This documentation
└── internal/
    └── calendar/
        └── google/
            └── client.go         # Google Calendar client implementation
```

## Additional Resources

- [Google Calendar API Documentation](https://developers.google.com/calendar/api)
- [Service Account Best Practices](https://cloud.google.com/iam/docs/service-accounts)
- [Google Cloud Console](https://console.cloud.google.com/)
- [Google Calendar](https://calendar.google.com)

## Support

If you encounter issues not covered in this documentation:

1. Check the troubleshooting section above
2. Verify your Google Cloud project settings
3. Ensure all prerequisites are met
4. Check the application logs for detailed error messages
5. Verify calendar sharing permissions

---

**Last Updated**: August 2025  
**Version**: 1.0  
**Maintainer**: Development Team

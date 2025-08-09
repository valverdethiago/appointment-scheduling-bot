# Quick Setup Guide

This is a condensed version of the setup process for developers who want to get started quickly.

## 🚀 Quick Start (5 minutes)

### 1. Prerequisites
- Google Cloud account
- Go development environment
- Access to Google Calendar

### 2. One-Time Setup

#### Create Service Account
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project or select existing
3. Enable Google Calendar API
4. Create service account: `appointment-scheduling-bot`
5. Download JSON key → rename to `service-account.json`
6. Place in project root

#### Configure Calendar
1. Go to [Google Calendar](https://calendar.google.com)
2. Share your calendar with service account email
3. Set permission: "Make changes to events"

#### Update Config
```bash
# Edit config.env
GCAL_CALENDAR_ID=your-calendar@gmail.com
GOOGLE_CREDS_JSON=service-account.json
```

### 3. Test Setup
```bash
# Quick test
./scripts/test-calendar.sh

# Full API test
go run scripts/test-calendar-connection.go
```

### 4. Start Development
```bash
# Start API server
go run cmd/api/main.go

# Health check
curl http://localhost:8080/healthz
```

## 📁 Key Files

- `service-account.json` - Google credentials (DO NOT COMMIT)
- `config.env` - Configuration
- `scripts/setup-calendar.sh` - Setup helper
- `scripts/test-calendar.sh` - Basic test
- `docs/google-calendar-setup.md` - Full documentation

## 🔧 Common Commands

```bash
# Test configuration
./scripts/test-calendar.sh

# Test API connection
go run scripts/test-calendar-connection.go

# Start API server
go run cmd/api/main.go

# Build binary
go build -o api cmd/api/main.go

# Run tests
go test ./...
```

## ⚠️ Important Notes

- **Never commit** `service-account.json`
- Service account needs calendar sharing permissions
- Test both read and write access
- Use dedicated calendar for production

## 🆘 Need Help?

- Check [full documentation](google-calendar-setup.md)
- Review troubleshooting section
- Verify calendar sharing settings
- Check service account permissions

---

**For detailed setup instructions, see**: [google-calendar-setup.md](google-calendar-setup.md)

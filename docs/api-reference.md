# API Reference

This document provides comprehensive documentation for all API endpoints, request/response formats, and usage examples for the Appointment Scheduling Bot.

## Table of Contents

- [Overview](#overview)
- [Base URL](#base-url)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Endpoints](#endpoints)
- [Data Models](#data-models)
- [Examples](#examples)
- [Rate Limiting](#rate-limiting)

## Overview

The Appointment Scheduling Bot API provides RESTful endpoints for managing appointments, checking availability, and handling scheduling operations. The API is built with Go and Fiber, providing fast and reliable performance.

## Base URL

- **Development**: `http://localhost:8080`
- **Production**: `https://your-domain.com`

## Authentication

Currently, the API uses service account authentication for Google Calendar operations. Future versions will include user authentication and authorization.

### Service Account Authentication

The API automatically authenticates with Google Calendar using the configured service account credentials in `service-account.json`.

## Error Handling

The API uses standard HTTP status codes and returns error details in JSON format:

```json
{
  "error": "Error message description",
  "code": "ERROR_CODE",
  "details": "Additional error details"
}
```

### Common HTTP Status Codes

- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (authentication required)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `500` - Internal Server Error (server error)

## Endpoints

### Health Check

#### GET `/healthz`

Returns the health status of the API and configuration information.

**Response:**
```json
{
  "status": "ok",
  "env": "development",
  "http_port": "8080",
  "timezone": "UTC",
  "gcal_calendar_id": "appointmentschedulingbot@gmail.com",
  "has_google_creds": true,
  "has_supabase_url": true,
  "has_supabase_key": true,
  "redis_url": "redis://localhost:6379"
}
```

### Calendar Operations

#### GET `/api/v1/calendar/slots`

Get available time slots for a given date range.

**Query Parameters:**
- `from` (required): Start date in ISO 8601 format (e.g., `2025-08-09T00:00:00Z`)
- `to` (required): End date in ISO 8601 format (e.g., `2025-08-09T23:59:59Z`)
- `duration` (optional): Appointment duration in minutes (default: 60)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/calendar/slots?from=2025-08-09T00:00:00Z&to=2025-08-09T23:59:59Z&duration=30"
```

**Response:**
```json
{
  "slots": [
    {
      "start": "2025-08-09T09:00:00Z",
      "end": "2025-08-09T09:30:00Z",
      "available": true
    },
    {
      "start": "2025-08-09T09:30:00Z",
      "end": "2025-08-09T10:00:00Z",
      "available": true
    }
  ],
  "total": 2,
  "date_range": {
    "from": "2025-08-09T00:00:00Z",
    "to": "2025-08-09T23:59:59Z"
  }
}
```

#### GET `/api/v1/calendar/busy`

Get busy time blocks for a given date range.

**Query Parameters:**
- `from` (required): Start date in ISO 8601 format
- `to` (required): End date in ISO 8601 format

**Response:**
```json
{
  "busy_times": [
    {
      "start": "2025-08-09T10:00:00Z",
      "end": "2025-08-09T11:00:00Z",
      "source": "google_calendar"
    }
  ],
  "total": 1
}
```

### Appointment Operations

#### POST `/api/v1/appointments`

Create a new appointment.

**Request Body:**
```json
{
  "summary": "Flu Shot Appointment",
  "description": "Annual flu vaccination",
  "start": "2025-08-09T09:00:00Z",
  "end": "2025-08-09T10:00:00Z",
  "attendee_name": "John Doe",
  "attendee_email": "john.doe@example.com",
  "location": "Main Pharmacy",
  "timezone": "America/New_York"
}
```

**Response:**
```json
{
  "id": "abc123def456",
  "summary": "Flu Shot Appointment",
  "description": "Annual flu vaccination",
  "start": "2025-08-09T09:00:00Z",
  "end": "2025-08-09T10:00:00Z",
  "attendee_name": "John Doe",
  "attendee_email": "john.doe@example.com",
  "location": "Main Pharmacy",
  "timezone": "America/New_York",
  "calendar_event_id": "google_calendar_event_id",
  "status": "confirmed",
  "created_at": "2025-08-09T08:00:00Z"
}
```

#### GET `/api/v1/appointments/{id}`

Get appointment details by ID.

**Response:**
```json
{
  "id": "abc123def456",
  "summary": "Flu Shot Appointment",
  "description": "Annual flu vaccination",
  "start": "2025-08-09T09:00:00Z",
  "end": "2025-08-09T10:00:00Z",
  "attendee_name": "John Doe",
  "attendee_email": "john.doe@example.com",
  "location": "Main Pharmacy",
  "timezone": "America/New_York",
  "calendar_event_id": "google_calendar_event_id",
  "status": "confirmed",
  "created_at": "2025-08-09T08:00:00Z",
  "updated_at": "2025-08-09T08:00:00Z"
}
```

#### PUT `/api/v1/appointments/{id}`

Update an existing appointment.

**Request Body:**
```json
{
  "summary": "Flu Shot Appointment - Updated",
  "start": "2025-08-09T10:00:00Z",
  "end": "2025-08-09T11:00:00Z"
}
```

**Response:**
```json
{
  "id": "abc123def456",
  "summary": "Flu Shot Appointment - Updated",
  "description": "Annual flu vaccination",
  "start": "2025-08-09T10:00:00Z",
  "end": "2025-08-09T11:00:00Z",
  "attendee_name": "John Doe",
  "attendee_email": "john.doe@example.com",
  "location": "Main Pharmacy",
  "timezone": "America/New_York",
  "calendar_event_id": "google_calendar_event_id",
  "status": "confirmed",
  "updated_at": "2025-08-09T08:30:00Z"
}
```

#### DELETE `/api/v1/appointments/{id}`

Cancel/delete an appointment.

**Response:**
```json
{
  "message": "Appointment cancelled successfully",
  "id": "abc123def456"
}
```

#### GET `/api/v1/appointments`

List appointments with optional filtering.

**Query Parameters:**
- `from` (optional): Start date filter
- `to` (optional): End date filter
- `attendee_email` (optional): Filter by attendee email
- `status` (optional): Filter by status (confirmed, cancelled, pending)
- `limit` (optional): Maximum results (default: 50)
- `offset` (optional): Pagination offset (default: 0)

**Response:**
```json
{
  "appointments": [
    {
      "id": "abc123def456",
      "summary": "Flu Shot Appointment",
      "start": "2025-08-09T09:00:00Z",
      "end": "2025-08-09T10:00:00Z",
      "attendee_name": "John Doe",
      "attendee_email": "john.doe@example.com",
      "status": "confirmed"
    }
  ],
  "total": 1,
  "limit": 50,
  "offset": 0
}
```

### Webhook Endpoints

#### POST `/api/v1/webhooks/retell`

Webhook endpoint for Retell.ai chat API integration.

**Request Body:**
```json
{
  "event": "tool_call",
  "data": {
    "tool_name": "find_slots",
    "parameters": {
      "date": "2025-08-09",
      "duration": 30
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "available_slots": [
      "2025-08-09T09:00:00Z",
      "2025-08-09T09:30:00Z"
    ]
  }
}
```

## Data Models

### Appointment

```json
{
  "id": "string",
  "summary": "string",
  "description": "string",
  "start": "datetime",
  "end": "datetime",
  "attendee_name": "string",
  "attendee_email": "string",
  "location": "string",
  "timezone": "string",
  "calendar_event_id": "string",
  "status": "confirmed|cancelled|pending",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### TimeSlot

```json
{
  "start": "datetime",
  "end": "datetime",
  "available": "boolean"
}
```

### TimeBlock

```json
{
  "start": "datetime",
  "end": "datetime",
  "source": "string"
}
```

## Examples

### Complete Appointment Booking Flow

1. **Check Availability:**
```bash
curl "http://localhost:8080/api/v1/calendar/slots?from=2025-08-09T00:00:00Z&to=2025-08-09T23:59:59Z&duration=60"
```

2. **Create Appointment:**
```bash
curl -X POST "http://localhost:8080/api/v1/appointments" \
  -H "Content-Type: application/json" \
  -d '{
    "summary": "Flu Shot",
    "description": "Annual vaccination",
    "start": "2025-08-09T09:00:00Z",
    "end": "2025-08-09T10:00:00Z",
    "attendee_name": "Jane Smith",
    "attendee_email": "jane.smith@example.com",
    "location": "Pharmacy",
    "timezone": "America/New_York"
  }'
```

3. **Confirm Appointment:**
```bash
curl "http://localhost:8080/api/v1/appointments/abc123def456"
```

### Error Handling Examples

**Invalid Date Format:**
```bash
curl "http://localhost:8080/api/v1/calendar/slots?from=invalid-date&to=2025-08-09T23:59:59Z"
```

**Response:**
```json
{
  "error": "Invalid date format",
  "code": "INVALID_DATE",
  "details": "Date parameter 'from' must be in ISO 8601 format"
}
```

**Appointment Not Found:**
```bash
curl "http://localhost:8080/api/v1/appointments/nonexistent-id"
```

**Response:**
```json
{
  "error": "Appointment not found",
  "code": "NOT_FOUND",
  "details": "No appointment found with ID: nonexistent-id"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default Limit**: 100 requests per minute per IP
- **Burst Limit**: 10 requests per second
- **Headers**: Rate limit information is included in response headers

### Rate Limit Headers

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## SDKs and Libraries

### Go Client

```go
package main

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
)

type AppointmentClient struct {
    baseURL string
    client  *http.Client
}

func (c *AppointmentClient) CreateAppointment(appt Appointment) (*Appointment, error) {
    jsonData, _ := json.Marshal(appt)
    resp, err := c.client.Post(c.baseURL+"/api/v1/appointments", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result Appointment
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

### JavaScript/Node.js Client

```javascript
class AppointmentClient {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }
    
    async createAppointment(appointment) {
        const response = await fetch(`${this.baseURL}/api/v1/appointments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(appointment),
        });
        
        return response.json();
    }
}
```

## Testing

### Test Endpoints

The API includes test endpoints for development and testing:

- **Health Check**: `/healthz`
- **Configuration**: `/api/v1/config` (development only)
- **Metrics**: `/metrics` (if enabled)

### Postman Collection

A Postman collection is available for testing all endpoints:

```json
{
  "info": {
    "name": "Appointment Scheduling Bot API",
    "description": "Complete API collection for testing"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": "{{base_url}}/healthz"
      }
    }
  ]
}
```

---

**Last Updated**: August 2025  
**Version**: 1.0  
**Maintainer**: Development Team

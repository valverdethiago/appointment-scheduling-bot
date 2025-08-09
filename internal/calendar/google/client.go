package google

import (
	"context"
	"fmt"
	"time"

	"appointment-scheduling-bot/internal/calendar"
	"appointment-scheduling-bot/internal/shared/config"

	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Client implements the calendar.Client interface for Google Calendar
type Client struct {
	service      *gcal.Service
	calendarID   string
}

// NewClient creates a new Google Calendar client
func NewClient(cfg config.Config) (*Client, error) {
	creds, err := cfg.GetGoogleCreds()
	if err != nil {
		return nil, fmt.Errorf("failed to get Google credentials: %w", err)
	}

	// Create JWT config from service account credentials
	jwtConfig, err := google.JWTConfigFromJSON(creds, gcal.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %w", err)
	}

	// Create OAuth2 client from JWT config
	client := jwtConfig.Client(context.Background())

	// Create calendar service
	service, err := gcal.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	return &Client{
		service:    service,
		calendarID: cfg.GCalCalendarID,
	}, nil
}

// ListBusy returns all busy time blocks between from and to
func (c *Client) ListBusy(from, to time.Time) ([]calendar.TimeBlock, error) {

	// Create free/busy query
	query := &gcal.FreeBusyRequest{
		TimeMin: from.Format(time.RFC3339),
		TimeMax: to.Format(time.RFC3339),
		Items: []*gcal.FreeBusyRequestItem{
			{Id: c.calendarID},
		},
	}

	// Execute free/busy query
	call := c.service.Freebusy.Query(query)
	result, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to query free/busy: %w", err)
	}

	var timeBlocks []calendar.TimeBlock
	
	// Process busy times for the calendar
	if cal, exists := result.Calendars[c.calendarID]; exists {
		for _, busy := range cal.Busy {
			start, err := time.Parse(time.RFC3339, busy.Start)
			if err != nil {
				continue
			}
			end, err := time.Parse(time.RFC3339, busy.End)
			if err != nil {
				continue
			}

			timeBlocks = append(timeBlocks, calendar.TimeBlock{
				Start:  start,
				End:    end,
				Source: "google_calendar",
			})
		}
	}

	return timeBlocks, nil
}

// CreateEvent creates a new calendar event and returns the event ID
func (c *Client) CreateEvent(appt calendar.Appointment) (string, error) {

	// Create calendar event
	event := &gcal.Event{
		Summary:     appt.Summary,
		Description: appt.Description,
		Start: &gcal.EventDateTime{
			DateTime: appt.Start.Format(time.RFC3339),
			TimeZone: appt.Timezone,
		},
		End: &gcal.EventDateTime{
			DateTime: appt.End.Format(time.RFC3339),
			TimeZone: appt.Timezone,
		},
		Location: appt.Location,
	}

	// Add attendee if provided
	if appt.AttendeeEmail != "" {
		event.Attendees = []*gcal.EventAttendee{
			{
				Email: appt.AttendeeEmail,
			},
		}
	}

	// Insert the event
	call := c.service.Events.Insert(c.calendarID, event)
	createdEvent, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("failed to create calendar event: %w", err)
	}

	return createdEvent.Id, nil
}

// UpdateEvent updates an existing calendar event
func (c *Client) UpdateEvent(eventID string, appt calendar.Appointment) error {

	// First, get the existing event
	getCall := c.service.Events.Get(c.calendarID, eventID)
	existingEvent, err := getCall.Do()
	if err != nil {
		return fmt.Errorf("failed to get existing event: %w", err)
	}

	// Update the event fields
	existingEvent.Summary = appt.Summary
	existingEvent.Description = appt.Description
	existingEvent.Start = &gcal.EventDateTime{
		DateTime: appt.Start.Format(time.RFC3339),
		TimeZone: appt.Timezone,
	}
	existingEvent.End = &gcal.EventDateTime{
		DateTime: appt.End.Format(time.RFC3339),
		TimeZone: appt.Timezone,
	}
	existingEvent.Location = appt.Location

	// Update attendees if provided
	if appt.AttendeeEmail != "" {
		existingEvent.Attendees = []*gcal.EventAttendee{
			{
				Email: appt.AttendeeEmail,
			},
		}
	}

	// Update the event
	updateCall := c.service.Events.Update(c.calendarID, eventID, existingEvent)
	_, err = updateCall.Do()
	if err != nil {
		return fmt.Errorf("failed to update calendar event: %w", err)
	}

	return nil
}

// DeleteEvent deletes a calendar event
func (c *Client) DeleteEvent(eventID string) error {

	call := c.service.Events.Delete(c.calendarID, eventID)
	err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}

	return nil
} 
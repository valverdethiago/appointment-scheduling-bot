package calendar

import "time"

// TimeBlock represents a busy time block from Google Calendar
type TimeBlock struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Source string    `json:"source"` // e.g., "google_calendar", "manual_block"
}

// Appointment represents a calendar appointment/event
type Appointment struct {
	Summary       string    `json:"summary"`
	Description   string    `json:"description"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	AttendeeName  string    `json:"attendee_name"`
	AttendeeEmail string    `json:"attendee_email"`
	Location      string    `json:"location"`
	Timezone      string    `json:"timezone"`
}

// Client interface for Google Calendar operations
type Client interface {
	// ListBusy returns all busy time blocks between from and to
	ListBusy(from, to time.Time) ([]TimeBlock, error)
	
	// CreateEvent creates a new calendar event and returns the event ID
	CreateEvent(appt Appointment) (eventID string, err error)
	
	// UpdateEvent updates an existing calendar event
	UpdateEvent(eventID string, appt Appointment) error
	
	// DeleteEvent deletes a calendar event
	DeleteEvent(eventID string) error
} 
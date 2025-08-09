package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"appointment-scheduling-bot/internal/calendar"
	"appointment-scheduling-bot/internal/calendar/google"
	"appointment-scheduling-bot/internal/shared/config"

	"github.com/spf13/cobra"
)

var (
	cfg     config.Config
	client  *google.Client
	rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "Appointment Scheduling Bot CLI",
		Long:  `A CLI tool for testing the appointment scheduling bot components.`,
	}
)

func init() {
	// Load configuration
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v", err)
	}

	// Initialize Google Calendar client if possible
	if cfg.GCalCalendarID != "" && cfg.GoogleCredsJSON != "" {
		client, err = google.NewClient(cfg)
		if err != nil {
			log.Printf("Warning: Failed to initialize Google Calendar client: %v", err)
		}
	}
}

var listBusyCmd = &cobra.Command{
	Use:   "list-busy",
	Short: "List busy time blocks from Google Calendar",
	Long:  `List all busy time blocks between the specified date range.`,
	Run: func(cmd *cobra.Command, args []string) {
		if client == nil {
			fmt.Println("Error: Google Calendar client not initialized")
			os.Exit(1)
		}

		fromStr, _ := cmd.Flags().GetString("from")
		toStr, _ := cmd.Flags().GetString("to")

		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			fmt.Printf("Error parsing from date: %v\n", err)
			os.Exit(1)
		}

		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			fmt.Printf("Error parsing to date: %v\n", err)
			os.Exit(1)
		}

		// Set time to start/end of day
		from = from.Add(24 * time.Hour) // Start of day after
		to = to.Add(24 * time.Hour)     // End of day

		timeBlocks, err := client.ListBusy(from, to)
		if err != nil {
			fmt.Printf("Error listing busy blocks: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Busy time blocks from %s to %s:\n", fromStr, toStr)
		if len(timeBlocks) == 0 {
			fmt.Println("No busy blocks found.")
		} else {
			for i, block := range timeBlocks {
				fmt.Printf("%d. %s - %s (%s)\n", 
					i+1, 
					block.Start.Format("2006-01-02 15:04"), 
					block.End.Format("2006-01-02 15:04"),
					block.Source)
			}
		}
	},
}

var createEventCmd = &cobra.Command{
	Use:   "create-event",
	Short: "Create a new calendar event",
	Long:  `Create a new calendar event with the specified details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if client == nil {
			fmt.Println("Error: Google Calendar client not initialized")
			os.Exit(1)
		}

		summary, _ := cmd.Flags().GetString("summary")
		description, _ := cmd.Flags().GetString("description")
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		attendeeName, _ := cmd.Flags().GetString("attendee-name")
		attendeeEmail, _ := cmd.Flags().GetString("attendee-email")
		location, _ := cmd.Flags().GetString("location")
		timezone, _ := cmd.Flags().GetString("timezone")

		if summary == "" {
			fmt.Println("Error: summary is required")
			os.Exit(1)
		}

		if startStr == "" || endStr == "" {
			fmt.Println("Error: start and end times are required")
			os.Exit(1)
		}

		// Parse times
		start, err := time.Parse("2006-01-02T15:04:05", startStr)
		if err != nil {
			fmt.Printf("Error parsing start time: %v\n", err)
			os.Exit(1)
		}

		end, err := time.Parse("2006-01-02T15:04:05", endStr)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			os.Exit(1)
		}

		appt := calendar.Appointment{
			Summary:       summary,
			Description:   description,
			Start:         start,
			End:           end,
			AttendeeName:  attendeeName,
			AttendeeEmail: attendeeEmail,
			Location:      location,
			Timezone:      timezone,
		}

		eventID, err := client.CreateEvent(appt)
		if err != nil {
			fmt.Printf("Error creating event: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Event created successfully with ID: %s\n", eventID)
	},
}

func main() {
	// Add flags to list-busy command
	listBusyCmd.Flags().String("from", time.Now().Format("2006-01-02"), "Start date (YYYY-MM-DD)")
	listBusyCmd.Flags().String("to", time.Now().AddDate(0, 0, 7).Format("2006-01-02"), "End date (YYYY-MM-DD)")

	// Add flags to create-event command
	createEventCmd.Flags().String("summary", "", "Event summary (required)")
	createEventCmd.Flags().String("description", "", "Event description")
	createEventCmd.Flags().String("start", "", "Start time (YYYY-MM-DDTHH:MM:SS) (required)")
	createEventCmd.Flags().String("end", "", "End time (YYYY-MM-DDTHH:MM:SS) (required)")
	createEventCmd.Flags().String("attendee-name", "", "Attendee name")
	createEventCmd.Flags().String("attendee-email", "", "Attendee email")
	createEventCmd.Flags().String("location", "", "Event location")
	createEventCmd.Flags().String("timezone", "UTC", "Event timezone")

	// Add commands to root
	rootCmd.AddCommand(listBusyCmd)
	rootCmd.AddCommand(createEventCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 
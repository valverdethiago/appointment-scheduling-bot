package main

import (
	"fmt"
	"log"
	"time"

	"appointment-scheduling-bot/internal/calendar"
	"appointment-scheduling-bot/internal/calendar/google"
	"appointment-scheduling-bot/internal/shared/config"
)

func main() {
	fmt.Println("🧪 Testing Google Calendar API Connection")
	fmt.Println("========================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	fmt.Printf("✅ Configuration loaded\n")
	fmt.Printf("   Calendar ID: %s\n", cfg.GCalCalendarID)
	fmt.Printf("   Service Account: %s\n", "appointment-scheduling-bot@rugged-precept-468518-a2.iam.gserviceaccount.com")

	// Create Google Calendar client
	fmt.Println("\n🔌 Creating Google Calendar client...")
	client, err := google.NewClient(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}
	fmt.Println("✅ Client created successfully")

	// Test calendar access by listing busy times
	fmt.Println("\n📅 Testing calendar access...")
	now := time.Now()
	from := now.Add(-24 * time.Hour)  // 24 hours ago
	to := now.Add(24 * time.Hour)     // 24 hours from now

	busyTimes, err := client.ListBusy(from, to)
	if err != nil {
		log.Fatalf("❌ Failed to list busy times: %v", err)
	}

	fmt.Printf("✅ Successfully accessed calendar!\n")
	fmt.Printf("   Found %d busy time blocks\n", len(busyTimes))

	// Test creating a test event (optional - you can comment this out if you don't want test events)
	fmt.Println("\n📝 Testing event creation (optional)...")
	
	// Create a test appointment
	testAppt := calendar.Appointment{
		Summary:        "Test Appointment - Bot Connection",
		Description:    "This is a test event to verify the service account connection",
		Start:         now.Add(1 * time.Hour),
		End:           now.Add(2 * time.Hour),
		Timezone:      "UTC",
		Location:      "Test Location",
		AttendeeEmail: "",
	}

	eventID, err := client.CreateEvent(testAppt)
	if err != nil {
		fmt.Printf("⚠️  Event creation failed (this might be expected): %v\n", err)
		fmt.Println("   This could mean the service account doesn't have write permissions yet")
	} else {
		fmt.Printf("✅ Test event created successfully!\n")
		fmt.Printf("   Event ID: %s\n", eventID)
		
		// Clean up the test event
		fmt.Println("🧹 Cleaning up test event...")
		if err := client.DeleteEvent(eventID); err != nil {
			fmt.Printf("⚠️  Failed to delete test event: %v\n", err)
		} else {
			fmt.Println("✅ Test event cleaned up")
		}
	}

	fmt.Println("\n🎉 Google Calendar Service Account Test Complete!")
	fmt.Println("================================================")
	fmt.Println("✅ Configuration: Working")
	fmt.Println("✅ Authentication: Working")
	fmt.Println("✅ Calendar Access: Working")
	fmt.Println("✅ API Connection: Working")
	
	if err != nil {
		fmt.Println("⚠️  Write Permissions: May need calendar sharing setup")
	} else {
		fmt.Println("✅ Write Permissions: Working")
	}
}

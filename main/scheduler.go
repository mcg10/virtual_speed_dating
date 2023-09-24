package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func scheduleMeetings(ctx context.Context, client *http.Client) error {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Calendar client: %v", err)
	}
	event := createEvent()
	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).ConferenceDataVersion(1).Do()
	if err != nil {
		return fmt.Errorf("Unable to create event. %v\n", err)
	}
	return nil
}

func createEvent() *calendar.Event {
	return &calendar.Event{
		Summary:     "Google I/O 2015",
		Location:    "800 Howard St., San Francisco, CA 94103",
		Description: "A chance to hear more about Google's developer products.",
		Start: &calendar.EventDateTime{
			DateTime: "2023-09-28T09:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: "2023-09-28T17:00:00-07:00",
			TimeZone: "America/Los_Angeles",
		},
		Attendees: []*calendar.EventAttendee{
			{Email: "matthew.giglio10@gmail.com"},
			{Email: "matthew@joincandidhealth.com"},
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: uuid.New().String(),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}
}

package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func scheduleMeetings(ctx context.Context, client *http.Client, day string, startTime string, duration int) error {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Calendar client: %v", err)
	}
	records, err := readFile()
	if err != nil {
		return err
	}
	for i := 0; i < len(records); i++ {
		if i != 0 {
			event := createEvent(records[i], day, startTime, duration)
			calendarId := "primary"
			event, err = srv.Events.Insert(calendarId, event).ConferenceDataVersion(1).Do()
			if err != nil {
				return fmt.Errorf("Unable to create event. %v\n", err)
			}
		}
	}
	return nil
}

func readFile() ([][]string, error) {
	file, err := os.Open("meetings.csv")
	if err != nil {
		return nil, fmt.Errorf("Error while reading the file %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading records: %v", err)
	}
	return records, nil
}

func getTimes(round, startTime string, duration int) (string, string) {
	round_num, err := strconv.Atoi(round)
	if err != nil {
		log.Fatal(err.Error())
	}
	t, err := time.Parse("15:04:05", startTime)
	if err != nil {
		log.Fatal(err.Error())
	}
	time_to_add := time.Duration(duration*(round_num-1))*time.Minute + time.Duration(round_num-1)*time.Minute
	start := t.Add(time_to_add)
	startingTime := start.Format("15:04:05")
	time_to_add = time.Duration(duration) * time.Minute
	end := start.Add(time_to_add)
	endingTime := end.Format("15:04:05")
	return startingTime, endingTime
}

func createEvent(record []string, day, startTime string, duration int) *calendar.Event {
	round, attendee_1, attendee_2 := record[0], record[1], record[2]
	startingTime, endingTime := getTimes(round, startTime, duration)
	return &calendar.Event{
		Summary: fmt.Sprintf("FSD Round %s", round),
		Start: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s", day, startingTime),
			TimeZone: "America/Los_Angeles",
		},
		End: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s", day, endingTime),
			TimeZone: "America/Los_Angeles",
		},
		Attendees: []*calendar.EventAttendee{
			{Email: attendee_1},
			{Email: attendee_2},
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

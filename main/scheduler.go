package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var timeFormat string = "15:04:05"
var calendarId string = "primary"
var timeZone string = "America/Los_Angeles"
var googleMeet string = "hangoutsMeet"

func scheduleMeetings(ctx context.Context, client *http.Client, day string, startTime string, duration int) error {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("Unable to retrieve Calendar client: %v", err)
	}
	records, err := readFile()
	if err != nil {
		return err
	}
	requests, err := createRequests(records, day, startTime, duration)
	if err != nil {
		return err
	}
	printRequests(requests)
	fmt.Println("If you wish to proceed, enter Yes. Otherwise, enter anything else")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("An error occured while reading input: %v", err)
	}
	input = strings.TrimSuffix(input, "\n")
	if input != "Yes" {
		return nil
	}
	for _, request := range requests {
		_, err = srv.Events.Insert(calendarId, request).ConferenceDataVersion(1).Do()
		if err != nil {
			return fmt.Errorf("Unable to create event. %v\n", err)
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

func getTimes(round, startTime string, duration int) (string, string, error) {
	round_num, err := strconv.Atoi(round)
	if err != nil {
		return "", "", fmt.Errorf("An error occurred parsing the round number: %v", err.Error())
	}
	t, err := time.Parse(timeFormat, startTime)
	if err != nil {
		return "", "", fmt.Errorf("An error occurred parsing the start time: %v", err.Error())
	}
	t = t.Add(time.Duration(duration*(round_num-1))*time.Minute + time.Duration(round_num-1)*time.Minute)
	startingTime := t.Format(timeFormat)
	t = t.Add(time.Duration(duration) * time.Minute)
	endingTime := t.Format(timeFormat)
	return startingTime, endingTime, nil
}

func createRequests(records [][]string, day, startTime string, duration int) ([]*calendar.Event, error) {
	var requests []*calendar.Event
	for i := 0; i < len(records); i++ {
		if i != 0 {
			event, err := createEvent(records[i], day, startTime, duration)
			if err != nil {
				return nil, err
			}
			requests = append(requests, event)
		}
	}
	return requests, nil
}

func createEvent(record []string, day, startTime string, duration int) (*calendar.Event, error) {
	round, attendee_1, attendee_2 := record[0], record[1], record[2]
	startingTime, endingTime, err := getTimes(round, startTime, duration)
	if err != nil {
		return nil, err
	}
	return &calendar.Event{
		Summary: fmt.Sprintf("FSD Round %s", round),
		Start: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s", day, startingTime),
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: fmt.Sprintf("%sT%s", day, endingTime),
			TimeZone: timeZone,
		},
		Attendees: []*calendar.EventAttendee{
			{Email: attendee_1},
			{Email: attendee_2},
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: uuid.New().String(),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: googleMeet,
				},
			},
		},
	}, nil
}

func printRequests(requests []*calendar.Event) {
	fmt.Println("Meeting sessions scheduled:")
	for _, request := range requests {
		fmt.Printf("%s %s - %s %s <> %s\n", request.Summary, request.Start.DateTime, request.End.DateTime,
			request.Attendees[0].Email, request.Attendees[1].Email)
	}
}

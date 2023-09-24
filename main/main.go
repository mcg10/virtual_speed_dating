package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getInfo(message string) (string, error) {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return "", fmt.Errorf("An error occured while reading input: %v", err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input, nil
}

func main() {
	ctx := context.Background()
	client, err := getCalendarClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	day, err := getInfo("Enter day in YYYY-MM-DD format: ")
	if err != nil {
		log.Fatal(err.Error())
	}
	startTime, err := getInfo("Enter starting time in HH:MM:SS format in PT: ")
	if err != nil {
		log.Fatal(err.Error())
	}
	sessionLength, err := getInfo("Enter length of session in minutes: ")
	if err != nil {
		log.Fatal(err.Error())
	}
	length, err := strconv.Atoi(sessionLength)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = scheduleMeetings(ctx, client, day, startTime, length)
	if err != nil {
		log.Fatal(err.Error())
	}
}

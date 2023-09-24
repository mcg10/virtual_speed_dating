package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	client, err := getCalendarClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = scheduleMeetings(ctx, client)
	if err != nil {
		log.Fatal(err.Error())
	}
}

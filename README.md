# Virtual Speed Dating Scheduler

This package can create Google Calendar events with associated Google Meet links for virtual speed dating sessions. You will need a CSV file with requested pairings and round numbers. The CSV must be named `meetings.csv` and be of the following format:
| Round    | Attendee 1 | Attendee 2
| -------- | ------- | ------- |
| 1  | person1@gmail.com   | person2@gmail.com |
| 2 | person2@gmail.com     | person3@gmail.com |
| 2    | person1@gmail.com    | person4@gmail.com |

## How to Use
1. Follow the [GCP instructions](https://developers.google.com/calendar/api/quickstart/go) on how to enable the Calendar API, configure the OAuth consent screen, and authorize credentials for a desktop app. Save the resulting credentials JSON as `credentials.json` under the `main` package.
2. In the `main` package, run `go build` and then `./main`.
3. Clink the link that appears in your Terminal and enter the resulting token into your Terminal. 
4. Once authenticated, you should be able to continually run the app by running `./main`. There will be specific command line prompts to enter the day of your session, the starting time of the overall session, and the duration of each round.
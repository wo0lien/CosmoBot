package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	// log token expiration
	log.Printf("Token expiration: %v", tok.Expiry)
	// refresh token if expired
	if tok.Expiry.Before(time.Now()) {
		tok, err = refreshToken(tok, *config)
		if err != nil {
			log.Fatalf("Unable to refresh token: %v", err)
		}
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// RefreshToken if expired using the config and the old token
func refreshToken(oldToken *oauth2.Token, config oauth2.Config) (*oauth2.Token, error) {
	// log
	log.Printf("Refreshing token...")
	newToken, err := config.TokenSource(context.TODO(), oldToken).Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v) %v\n", item.Summary, date, item.Id)
		}
	}

	// eventID := "<event_id>"
	// calendarID := "primary"

	// // Retrieve the event
	// event, err := srv.Events.Get(calendarID, eventID).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve event: %v", err)
	// }

	// // Add an attendee to the event
	// newAttendee := &calendar.EventAttendee{
	// 	Email: "antoine.merle@insa-lyon.fr",
	// }
	// event.Attendees = append(event.Attendees, newAttendee)

	// event.Visibility = "default"

	// // Update the event
	// updatedEvent, err := srv.Events.Update(calendarID, eventID, event).SendUpdates("all").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to update event: %v", err)
	// }

	// fmt.Printf("Event updated: %s\n", updatedEvent.Id)
}

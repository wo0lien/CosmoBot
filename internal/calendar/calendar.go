package calendar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var service Service

type Service struct {
	*calendar.Service
}

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

// Starts the calendar service
func CalendarService() (*Service, error) {
	if service.Service != nil {
		return &service, nil
	}

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, err
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &Service{srv}, err
}

// Get event in the calendar
func (s *Service) Event(event *models.CosmoEvent) (*calendar.Event, error) {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return nil, errors.New("calendar event does not exist")
	}
	ev, err := s.Events.Get("primary", *event.CalendarID).Do()
	if err != nil {
		return nil, err
	}
	return ev, nil
}

// Create a new event in the calendar
// Does not update the DB with the google eventID of the event
func (s *Service) CreateEvent(event *models.CosmoEvent) (*calendar.Event, error) {
	// create event
	newEvent := &calendar.Event{
		Summary: event.Name,
		// TODO
		Location:    "",
		Description: "",
		Start: &calendar.EventDateTime{
			DateTime: event.StartDate.Format(time.RFC3339),
			TimeZone: "Europe/Paris",
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndDate.Format(time.RFC3339),
			TimeZone: "Europe/Paris",
		},
	}

	// insert event
	ev, err := s.Events.Insert("primary", newEvent).Do()
	if err != nil {
		return nil, err
	}

	return ev, nil
}

// Update event in the calendar
func (s *Service) UpdateEvent(event *models.CosmoEvent) (*calendar.Event, error) {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return nil, errors.New("calendar event does not exist")
	}

	// retrieve event
	ev, err := s.Events.Get("primary", *event.CalendarID).Do()
	if err != nil {
		return nil, err
	}

	// update event
	ev.Summary = event.Name
	ev.Start.DateTime = event.StartDate.Format(time.RFC3339)
	ev.End.DateTime = event.EndDate.Format(time.RFC3339)

	// update event
	updatedEvent, err := s.Events.Update("primary", *event.CalendarID, ev).SendUpdates("all").Do()
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// Update attendees of an event in the calendar
// Do not update the event in the database
// Replace all attendees of the event by the new ones
func (s *Service) UpdateEventAttendees(event *models.CosmoEvent, volunteers *[]models.Volunteer) (*calendar.Event, error) {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return nil, errors.New("calendar event does not exist")
	}

	// retrieve event
	ev, err := s.Events.Get("primary", *event.CalendarID).Do()
	if err != nil {
		return nil, err
	}

	// remove all attendees
	ev.Attendees = []*calendar.EventAttendee{}

	// update attendees
	for _, volunteer := range *volunteers {
		// add volunteer to attendees
		newAttendee := &calendar.EventAttendee{
			Email: volunteer.Email,
		}
		ev.Attendees = append(ev.Attendees, newAttendee)
	}

	// update event
	updatedEvent, err := s.Events.Update("primary", *event.CalendarID, ev).SendUpdates("all").Do()
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (s *Service) AddEventAttendee(event *models.CosmoEvent, volunteer *models.Volunteer) (*calendar.Event, error) {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return nil, errors.New("calendar event does not exist")
	}

	// retrieve event
	ev, err := s.Events.Get("primary", *event.CalendarID).Do()
	if err != nil {
		return nil, err
	}

	// append attendee to event list
	ev.Attendees = append(ev.Attendees, &calendar.EventAttendee{
		Email: volunteer.Email,
	})

	// update event
	updatedEvent, err := s.Events.Update("primary", *event.CalendarID, ev).SendUpdates("all").Do()
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (s *Service) RemoveEventAttendee(event models.CosmoEvent, volunteer models.Volunteer) (*calendar.Event, error) {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return nil, errors.New("calendar event does not exist")
	}

	// retrieve event
	ev, err := s.Events.Get("primary", *event.CalendarID).Do()
	if err != nil {
		return nil, err
	}

	// remove attendee from event list
	for i, attendee := range ev.Attendees {
		if attendee.Email == volunteer.Email {
			ev.Attendees = append(ev.Attendees[:i], ev.Attendees[i+1:]...)
			break
		}
	}

	// update event
	updatedEvent, err := s.Events.Update("primary", *event.CalendarID, ev).SendUpdates("all").Do()
	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

// Delete event in the calendar
// Do not delete the event in the database neither the calendarID of the event or flag
func (s *Service) DeleteEvent(event *models.CosmoEvent) error {
	// checks if calendar exists
	if !event.DoesCalendarExist {
		return errors.New("calendar event does not exist")
	}

	// delete event
	err := s.Events.Delete("primary", *event.CalendarID).Do()
	if err != nil {
		return err
	}

	return nil
}

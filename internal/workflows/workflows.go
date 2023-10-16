package workflows

import (
	"fmt"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/calendar"
	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"github.com/wo0lien/cosmoBot/internal/utils"
)

// Modules aims to be the main logic parts of the program
/*
TODO Send reminders of communication if it does not exist at desired timing
TODO Send rex reminders
TODO Update calendar on even update
*/

func StartDiscussionForEvent(event *models.CosmoEvent) error {
	if event.DoesChannelExist {
		logging.Info.Printf("Event %s already has a channel, skipping it", event.Name)
		return nil
	}
	ch, err := discord.Bot.StartEventDiscussion(event, fmt.Sprintf("%s - %s", event.StartDate.Format("01/02"), event.Name), ":Cosmix:")

	if err != nil {
		return err
	}

	event.DoesChannelExist = true
	event.ChannelID = &ch.ID
	err = controllers.SaveEvent(event)
	if err != nil {
		return err
	}
	return nil
}

// Lookup for upcoming events in DB and create a discussion for each of them
// if the discussion does not exist
func StartDiscussionForUpcomingEvents() error {
	logging.Info.Println("Starting discussion for upcoming events")
	// use events from db
	events := controllers.AllUpcomingEvents()

	for _, ev := range *events {
		err := StartDiscussionForEvent(&ev)
		if err != nil {
			logging.Error.Println(err)
		}
	}

	return nil
}

func TagAllVolunteersInAllUpcomingEvents() {
	// get all volunteers events joins
	joins, err := controllers.AllUpcomingVolunteersEvents()
	if err != nil {
		logging.Error.Printf("got an error getting all volunteers events : %s", err)
		return
	}

	// for each join
	for _, join := range *joins {
		// check if already tagged
		if !join.HasBeenTagged {
			// get volunteer and event
			volunteer, err := controllers.VolunteerById(join.VolunteerID)
			if err != nil {
				logging.Error.Println(err)
				continue
			}

			event, err := controllers.EventByID(join.CosmoEventID)
			if err != nil {
				logging.Error.Println(err)
				continue
			}

			err = TagVolunteerInEvent(volunteer, event, &join)
			if err != nil {
				logging.Error.Println(err)
				continue
			}
		}
	}
}

// TagVolunteerInEvent tags a volunteer in an event
// It also updates the join in the database
// If the join is not provided (nil), it will be fetched from the database
func TagVolunteerInEvent(volunteer *models.Volunteer, event *models.CosmoEvent, join *models.VolunteerEvent) error {
	// check if volunteer has a discord ID
	if volunteer.DiscordID == nil {
		return fmt.Errorf("volunteer %s %s does not have a discordID set", volunteer.FirstName, volunteer.LastName)
	}

	// check if event has a channel
	if !event.DoesChannelExist {
		logging.Info.Printf("event %s does not have a channel", event.Name)
		return nil
	}

	if join == nil {
		// putting var err error here to use the same join variable
		var err error
		join, err = controllers.VolunteerEventJoinByVolunteerIDAndEventID(volunteer.ID, event.ID)
		if err != nil {
			return err
		}
	}

	// check if already tagged
	if join.HasBeenTagged {
		logging.Info.Printf("Volunteer %s has already been tagged in event %s", volunteer.FirstName, event.Name)
		return nil
	}

	// tag volunteer
	logging.Info.Printf("Tagging volunteer %s in event %s", volunteer.FirstName, event.Name)
	_, err := discord.Bot.ChannelMessageSend(*event.ChannelID, fmt.Sprintf("<@%s> rejoins l'orga !", *volunteer.DiscordID))

	if err != nil {
		return err
	}

	// set tagged to true
	join.HasBeenTagged = true
	// save join
	err = controllers.SaveVolunteerEvent(join)
	if err != nil {
		logging.Error.Println(err)
	}

	return nil
}

func InviteVolunteersToAllUpcomingEvents(cs *calendar.Service) {
	// get all volunteers events joins
	joins, err := controllers.AllUpcomingVolunteersEvents()
	if err != nil {
		logging.Error.Printf("got an error getting all volunteers events : %s", err)
		return
	}

	// for each join
	for _, join := range *joins {
		// check if already tagged
		if !join.HasBeenInvited {
			// get volunteer and event
			volunteer, err := controllers.VolunteerById(join.VolunteerID)
			if err != nil {
				logging.Error.Println(err)
				continue
			}

			event, err := controllers.EventByID(join.CosmoEventID)
			if err != nil {
				logging.Error.Println(err)
				continue
			}

			err = InviteVolunteerInEvent(cs, volunteer, event)
			if err != nil {
				logging.Error.Println(err)
				continue
			}
			// set tagged to true
			join.HasBeenInvited = true
			// save join
			err = controllers.SaveVolunteerEvent(&join)
			if err != nil {
				logging.Error.Println(err)
				continue
			}

		}
	}
}

func InviteVolunteerInEvent(service *calendar.Service, volunteer *models.Volunteer, event *models.CosmoEvent) error {
	// check if event has a channel
	if !event.DoesCalendarExist {
		logging.Info.Printf("event %s does not have a calendar event", event.Name)
		return nil
	}

	logging.Info.Printf("Adding volunteer %s in event %s", volunteer.FirstName, event.Name)
	_, err := service.AddEventAttendee(event, volunteer)

	if err != nil {
		return err
	}

	return nil
}

func UninviteVolunteerInEvent(service *calendar.Service, volunteer models.Volunteer, event models.CosmoEvent) error {

	// check if event has a channel
	if !event.DoesCalendarExist {
		logging.Info.Printf("event %s does not have a calendar event", event.Name)
		return nil
	}

	logging.Info.Printf("Removing volunteer %s in event %s", volunteer.FirstName, event.Name)
	_, err := service.RemoveEventAttendee(event, volunteer)

	if err != nil {
		return err
	}

	return nil
}

// Create a calendar event and save it in the database
// Checks if the event has a calendar event first
func CreateCalendarEvent(cs *calendar.Service, event *models.CosmoEvent) error {
	// check if event has a channel
	if event.DoesCalendarExist {
		logging.Info.Printf("event %s already has a calendar event. Skipping it", event.Name)
		return nil
	}

	logging.Info.Printf("Creating calendar event for event %s", event.Name)
	calendarEvent, err := cs.CreateEvent(event)

	if err != nil {
		return err
	}

	event.CalendarID = &calendarEvent.Id
	event.DoesCalendarExist = true

	err = controllers.SaveEvent(event)
	if err != nil {
		return err
	}

	return nil
}

// Create a calendar event for each upcoming event
// update the database
func CrateCalendarEventForAllUpcomingEvents(cs *calendar.Service) {
	logging.Info.Println("Creating calendar events for all upcoming events")
	events := controllers.AllUpcomingEvents()

	for _, ev := range *events {
		err := CreateCalendarEvent(cs, &ev)
		if err != nil {
			logging.Error.Println(err)
		}
	}
}

// Adding a new relation from volunteer to event COULD trigger a volunteer insert
func InsertVolunteerByID(cs *calendar.Service, volunteerID uint) error {
	logging.Info.Println("Handling volunteer insert")

	volunteer, err := api.NocoApi.VolunteerByID(volunteerID)

	if err != nil {
		return err
	}

	// TODO use addedEventsIDs list to update calendar and discord
	volunteerInDB, _, _, err := controllers.UpdateVolunteerEventsJoinsFromVolunteerInApi(volunteer)

	// discord and calendar
	for _, event := range volunteerInDB.Events {
		// check if event is upcoming
		if !event.IsUpcoming() {
			continue
		}
		err := TagVolunteerInEvent(volunteerInDB, &event, nil)
		if err != nil {
			// just log the error, do not stop the process
			logging.Error.Println(err)
		}

		err = InviteVolunteerInEvent(cs, volunteerInDB, &event)
		if err != nil {
			// just log the error, do not stop the process
			logging.Error.Println(err)
		}
	}

	return err
}

// Removing a relation from volunteer to event COULD trigger a volunteer update
// Changing informations of a volunteer triggers a volunteer update
func UpdateVolunteerByID(cs *calendar.Service, volunteerID uint) error {
	logging.Info.Println("Handling volunteer update")

	volunteer, err := api.NocoApi.VolunteerByID(volunteerID)
	if err != nil {
		return err
	}

	volunteerInDB, err := controllers.CreateOrUpdateVolunteerToDBFromAPI(volunteer)
	if err != nil {
		return err
	}

	_, _, removedEventsIDs, err := controllers.UpdateVolunteerEventsJoinsFromVolunteerInApi(volunteer)
	if err != nil {
		return err
	}

	// calendar update (no discord update required)
	for _, eventID := range removedEventsIDs {
		event, err := controllers.EventByID(eventID)
		if err != nil {
			logging.Error.Println(err)
			continue
		}
		if event.IsUpcoming() {
			err = UninviteVolunteerInEvent(cs, *volunteerInDB, *event)
			if err != nil {
				logging.Error.Println(err)
			}
		}
	}

	return nil
}

// Removing a volunteer triggers a volunteer delete
func DeleteVolunteerByID(volunteerID uint) error {
	logging.Info.Println("Handling volunteer delete")

	_, err := controllers.VolunteerById(volunteerID)
	if err != nil {
		return err
	}

	err = controllers.DeleteVolunteerById(volunteerID)
	if err != nil {
		return err
	}

	return nil
}

// Adding a new relation from volunteer to event COULD trigger an event insert
// Creating a new event triggers an event insert
func InsertEventByID(cs *calendar.Service, eventID uint) error {
	logging.Info.Println("Handling event insert")

	// get event from api
	event, err := api.NocoApi.EventByID(eventID)
	if err != nil {
		return err
	}

	var eventInDB *models.CosmoEvent

	// get event from db
	eventInDB, err = controllers.EventByID(eventID)
	if err != nil {
		// event does not exist in db
		// create it
		eventInDB, err = controllers.CreateOrUpdateEventInDBFromApi(*event)
		if err != nil {
			return err
		}
	}

	// update volunteersEvent joins
	_, _, _, err = controllers.UpdateVolunteerEventsJoinsFromEventInApi(event)
	if err != nil {
		return err
	}

	// discord and calendar

	// check if event is upcoming
	if eventInDB.IsUpcoming() {
		// Start discussion for event
		_, err = discord.Bot.StartEventDiscussion(eventInDB, fmt.Sprintf("%s - %s", eventInDB.StartDate.Format("01/02"), eventInDB.Name), ":Cosmix:")
		if err != nil {
			logging.Error.Println(err)
		}

		// create calendar event
		err = CreateCalendarEvent(cs, eventInDB)
		if err != nil {
			logging.Error.Println(err)
		}

		// get all event’s volunteers
		volunteers, err := controllers.AllVolunteersByEventID(eventInDB.ID)
		if err != nil {
			return err
		}

		// tag all volunteers
		for _, volunteer := range *volunteers {
			err = TagVolunteerInEvent(&volunteer, eventInDB, nil)
			if err != nil {
				logging.Error.Println(err)
			}

			err = InviteVolunteerInEvent(cs, &volunteer, eventInDB)
			if err != nil {
				logging.Error.Println(err)
			}

		}
	}

	return nil
}

// Removing a relation from volunteer to event COULD trigger an event update
// Changing informations of an event triggers a volunteer update
func UpdateEventByID(cs *calendar.Service, eventID uint) error {
	logging.Info.Println("Handling event update")

	// get event from api
	event, err := api.NocoApi.EventByID(eventID)
	if err != nil {
		return err
	}

	var eventInDB *models.CosmoEvent

	//TODO improve to do everything in a single Api call
	// get event from db
	eventInDB, err = controllers.EventByID(eventID)
	if err != nil {
		return err
	}
	if eventInDB == nil {
		logging.Warning.Printf("event %d does not exist in db, creating it", eventID)
		// create it
		_, err = controllers.CreateOrUpdateEventInDBFromApi(*event)
		if err != nil {
			return err
		}
	}

	// update volunteersEvent joins
	eventInDB, _, removedEventsIDs, err := controllers.UpdateVolunteerEventsJoinsFromEventInApi(event)
	if err != nil {
		return err
	}

	if eventInDB.IsUpcoming() {
		// remove all volunteers that have been removed of the event
		volunteers, err := controllers.AllVolunteersByEventID(eventInDB.ID)
		if err != nil {
			return err
		}
		for _, volunteer := range *volunteers {
			if !utils.Contains(removedEventsIDs, volunteer.ID) {
				err = UninviteVolunteerInEvent(cs, volunteer, *eventInDB)
				if err != nil {
					logging.Error.Println(err)
				}
			}
		}
	}

	return nil
}

// Removing an event triggers an event delete
func DeleteEventByID(cs *calendar.Service, eventID uint) error {
	logging.Info.Println("Handling event delete")

	// get event from db
	eventInDB, err := controllers.EventWithVolunteersByID(eventID)
	if err != nil {
		return err
	}

	err = cs.DeleteEvent(eventInDB)
	if err != nil {
		logging.Error.Println(err)
	}

	// delete event from db
	err = controllers.DeleteEventById(eventID)
	if err != nil {
		return err
	}

	return nil
}

// Refresh all
func RefreshAll(cs *calendar.Service) {

	// load volunteers
	volunteers, err := api.NocoApi.AllVolunteers()
	if err != nil {
		panic(err)
	}

	// load events
	events, err := api.NocoApi.AllEvents()
	if err != nil {
		panic(err)
	}

	// put in db
	err = controllers.CreateOrUpdateEventsInDBFromApi(*events)
	if err != nil {
		panic(err)
	}

	err = controllers.CreateOrUpdateVolunteersToDBFromAPI(volunteers)
	if err != nil {
		panic(err)
	}

	err = controllers.UpdateVolunteersEventsJoinsFromApi(volunteers)
	if err != nil {
		panic(err)
	}

	// discord init
	StartDiscussionForUpcomingEvents()
	TagAllVolunteersInAllUpcomingEvents()

	// calendar init
	CrateCalendarEventForAllUpcomingEvents(cs)
	InviteVolunteersToAllUpcomingEvents(cs)
}

package modules

import (
	"fmt"

	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
)

// Modules aims to be the main logic parts of the program
/*

5. Send reminders of communication if it does not exist at desired timing
*/

// Lookup for upcoming events in DB and create a discussion for each of them
// if the discussion does not exist
func StartDiscussionForUpcomingEvents() error {
	logging.Info.Println("Starting discussion for upcoming events")
	// use events from db
	events := controllers.AllUpcomingEvents()

	for _, ev := range *events {
		logging.Info.Printf("Checking if event %s (id %d) has a channel", ev.Name, ev.ID)
		if !ev.DoesChannelExist {
			logging.Info.Printf("Event %s does not have a channel, creating one", ev.Name)
			ch, err := discord.Bot.StartEventDiscussion(&ev, fmt.Sprintf("%s - %s", ev.StartDate.Format("01/02"), ev.Name), ":Cosmix:")

			if err != nil {
				logging.Error.Printf("Could not create channel for event %s. Error: %s", ev.Name, err)
				continue
			}

			ev.DoesChannelExist = true
			ev.ChannelID = &ch.ID
			controllers.SaveEvent(&ev)
		}
	}

	return nil
}

func tagVolunteerInEvent(volunteerId, eventId uint) error {
	event, err := controllers.EventByID(eventId)
	if err != nil {
		return err
	}

	volunteer := controllers.GetVolunteerById(volunteerId)

	if event == nil || volunteer == nil {
		return fmt.Errorf("could not find event or volunteer with id %d and %d", eventId, volunteerId)
	}

	// check if event has a channel
	if !event.DoesChannelExist {
		logging.Info.Printf("event %s does not have a channel", event.Name)
		return nil
	}

	logging.Info.Printf("Tagging volunteer %s in event %s", volunteer.FirstName, event.Name)
	_, err = discord.Bot.ChannelMessageSend(*event.ChannelID, fmt.Sprintf("<@%s> rejoins l'orga !", *volunteer.DiscordID))

	if err != nil {
		return err
	}

	return nil
}

func TagAllVolunteersInAllEvents() {
	// get all volunteers events joins
	joins := controllers.GetAllVolunteersEvents()

	// for each join
	for _, join := range *joins {
		// check if already tagged
		if !join.VolunteerHasBeenTagged {
			// tag
			err := tagVolunteerInEvent(join.VolunteerID, join.CosmoEventID)
			if err != nil {
				logging.Error.Println(err)
			}
			// set tagged to true
			join.VolunteerHasBeenTagged = true
			// save join
			err = controllers.SaveVolunteerEvent(&join)
			if err != nil {
				logging.Error.Println(err)
			}

		}
	}
}

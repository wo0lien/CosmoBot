package modules

import (
	"fmt"

	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
)

// Modules aims to be the main logic parts of the program
/*

1. Load every events from NocoDB, store them in database.
Using NocoDB IDs will prevent adding two times the same element
Detecting deleted elements is trickier should be done by storing a list of IDs and poping
element from it for each time db event

2. If el in future + discussionId empty = start discussion

3. Tag volunteers when assigned to an event
If a volunteer can’t be tagged - Send a message to an admin channel

5. Send reminders of communication if it does not exist at desired timing

6.


TODO
2
DOING

DONE
1
*/

// Lookup for upcoming events in DB and create a discussion for each of them
// if the discussion does not exist
func StartDiscussionForUpcomingEvents() error {
	// use events from db
	events := controllers.GetAllUpcomingCosmoEvents()

	for _, ev := range *events {
		if !ev.DoesChannelExist {
			ch, err := discord.Bot.StartEventDiscussion(&ev, fmt.Sprintf("%s - %s", ev.StartDate.Format("01/02"), ev.Name), ":Cosmix:")

			if err != nil {
				return err
			}

			ev.DoesChannelExist = true
			ev.ChannelID = &ch.ID
			controllers.SaveEvent(&ev)
		}
	}

	return nil
}

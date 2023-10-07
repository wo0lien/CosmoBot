package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wo0lien/cosmoBot/internal/api/webhooks"
	"github.com/wo0lien/cosmoBot/internal/calendar"
)

func main() {
	calendar.Main()
	// StartNoco()
}

func StartNoco() {

	go webhooks.StartWebHooksHandlingServer()

	// events, err := api.NocoApi.GetAllEvents()

	// if err != nil {
	// 	panic(err)
	// }
	// controllers.LoadEventsInDBFromAPI(*events)

	// modules.StartDiscussionForUpcomingEvents()

	// volunteers, err := api.NocoApi.GetAllVolunteers()

	// if err != nil {
	// 	panic(err)
	// }

	// controllers.LoadVolunteersToDBFromAPI(volunteers)
	// controllers.LoadVolunteersEventsJoinsFromApi(volunteers)

	// modules.TagAllVolunteersInAllEvents()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

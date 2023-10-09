package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/api/webhooks"
	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/modules"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
)

func main() {
	// calendar.Main()
	Start()
	StartNoco()
}

func StartNoco() {

	go webhooks.StartWebHooksHandlingServer()

	// loading bot
	var _ = discord.Bot

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func Start() {
	// load volunteers
	volunteers, err := api.NocoApi.GetAllVolunteers()
	if err != nil {
		panic(err)
	}

	// load events
	events, err := api.NocoApi.GetAllEvents()
	if err != nil {
		panic(err)
	}

	// put in db
	err = controllers.LoadEventsInDBFromAPI(*events)
	if err != nil {
		panic(err)
	}

	err = controllers.LoadVolunteersToDBFromAPI(volunteers)
	if err != nil {
		panic(err)
	}

	err = controllers.LoadVolunteersEventsJoinsFromApi(volunteers)
	if err != nil {
		panic(err)
	}

	modules.StartDiscussionForUpcomingEvents()
	// modules.TagAllVolunteersInAllEvents()

}

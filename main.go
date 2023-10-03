package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/modules"
	"github.com/wo0lien/cosmoBot/internal/storage/controllers"
)

func main() {
	StartNoco()
}

func StartNoco() {
	events, err := api.NocoApi.GetAllEvents()

	if err != nil {
		panic(err)
	}
	controllers.LoadEventsInDBFromAPI(*events)

	modules.StartDiscussionForUpcomingEvents()

	volunteers, err := api.NocoApi.GetAllVolunteers()

	if err != nil {
		panic(err)
	}

	controllers.LoadVolunteersToDBFromAPI(volunteers)
	controllers.LoadVolunteersEventsJoinsFromApi(volunteers)

	modules.TagAllVolunteersInAllEvents()

}

func StartBot() {

	logging.Info.Println("CosmoBot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Bot.Close()
}

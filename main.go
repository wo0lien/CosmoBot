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
	upcomingEvents, err := api.NocoApi.GetAllUpcomingEvents()

	if err != nil {
		panic(err)
	}
	controllers.LoadEventsInDBFromAPI(*upcomingEvents)

	modules.StartDiscussionForUpcomingEvents()
}

func StartBot() {

	logging.Info.Println("CosmoBot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Bot.Close()
}

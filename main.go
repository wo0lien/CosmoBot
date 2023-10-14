package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/wo0lien/cosmoBot/internal/api/webhooks"
	"github.com/wo0lien/cosmoBot/internal/calendar"
	"github.com/wo0lien/cosmoBot/internal/discord"
	"github.com/wo0lien/cosmoBot/internal/workflows"
)

func main() {
	// loading bot
	var _ = discord.Bot

	cs, err := calendar.CalendarService()
	if err != nil {
		panic(err)
	}

	// Refresh everything
	workflows.RefreshAll(cs)

	// Start the webhooks server
	go webhooks.StartWebHooksHandlingServer()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

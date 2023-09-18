package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wo0lien/cosmoBot/internal/discord"
)

func main() {
	fmt.Println("CosmoBot is starting...")

	dg, err := discord.CreateBot()

	if err != nil {
		panic(-1)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

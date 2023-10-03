package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wo0lien/cosmoBot/internal/logging"
)

var Bot *discordBot

type discordBot struct {
	*discordgo.Session
}

var BOT_TOKEN string

func init() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	// load env variables BOT_TOKEN
	BOT_TOKEN = os.Getenv("BOT_TOKEN")
	if BOT_TOKEN == "" {
		logging.Critical.Fatal("BOT_TOKEN env variable is not set")
	}
	// start bot

	logging.Info.Println("CosmoBot is starting...")

	bot, err := discordgo.New("Bot " + BOT_TOKEN)

	if err != nil {
		logging.Critical.Fatalf("Could not start bot: %s", err)
	}

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	bot.AddHandler(messageCreate)

	// running the bot
	err = bot.Open()
	if err != nil {
		logging.Critical.Fatalf("Could not connect to discord: %s", err)
	}

	Bot = &discordBot{bot}

}

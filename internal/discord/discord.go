package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

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
		panic("BOT_TOKEN env variable is not set")
	}
}

// Instantiate a new bot
func CreateBot() (*discordgo.Session, error) {
	dg, err := discordgo.New("Bot " + BOT_TOKEN)

	if err != nil {
		return nil, err
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// running the bot

	err = dg.Open()
	if err != nil {
		return nil, err
	}

	return dg, nil
}

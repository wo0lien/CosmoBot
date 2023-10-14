package discord

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wo0lien/cosmoBot/internal/logging"
)

var Bot *discordBot

type discordBot struct {
	*discordgo.Session
	registreredCommands []*discordgo.ApplicationCommand
}

var BOT_TOKEN string
var DISCORD_GUILD_ID string

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

	// load env variables BOT_TOKEN
	DISCORD_GUILD_ID = os.Getenv("DISCORD_GUILD_ID")
	if DISCORD_GUILD_ID == "" {
		logging.Critical.Fatal("DISCORD_GUILD_ID env variable is not set")
	}
	// start bot

	logging.Info.Println("CosmoBot is starting...")

	bot, err := discordgo.New("Bot " + BOT_TOKEN)

	if err != nil {
		logging.Critical.Fatalf("Could not start bot: %s", err)
	}

	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// running the bot
	err = bot.Open()
	if err != nil {
		logging.Critical.Fatalf("Could not connect to discord: %s", err)
	}

	// adding handlers
	bot.AddHandler(messageCreate)
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	logging.Debug.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		logging.Debug.Printf("Registering command '%v'\n", v.Name)
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, DISCORD_GUILD_ID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	Bot = &discordBot{
		bot,
		registeredCommands,
	}

}

// Delete channel of the given id
func (d *discordBot) DeleteChannel(channelId string) error {
	_, err := d.ChannelDelete(channelId)
	if err != nil {
		return err
	}
	return nil
}

func (d *discordBot) Close() {
	d.Session.Close()
	for _, v := range d.registreredCommands {
		err := d.ApplicationCommandDelete(d.State.User.ID, DISCORD_GUILD_ID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

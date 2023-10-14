package discord

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/wo0lien/cosmoBot/internal/logging"
)

var NOCO_NEW_VOLUNTEER_FORM_URL string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	// load env variables BOT_TOKEN
	NOCO_NEW_VOLUNTEER_FORM_URL = os.Getenv("NOCO_NEW_VOLUNTEER_FORM_URL")
	if NOCO_NEW_VOLUNTEER_FORM_URL == "" {
		panic("NOCO_NEW_VOLUNTEER_FORM_URL env variable is not set")
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "bonjour",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Dis bonjour à Cosmix pour qu’il te connaisse !",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bonjour": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			logging.Info.Println("Received bonjour command")
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Content is a temporary field that will be replaced by Embeds
					// or files.
					Content: fmt.Sprintf("Salut ! Je te connais pas encore, tu peux remplir tes informations ici : %s (quand on te le demandera, ton id discord c’est %s)", NOCO_NEW_VOLUNTEER_FORM_URL, i.Member.User.ID),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		},
	}
)

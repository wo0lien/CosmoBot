package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
)

// starts either a post or a public thread depending on channelId with name and content
// post if the type is ChannelTypeGuildForum
// thread if the type is ChannelTypeGuildText
func (d *discordBot) StartDiscussion(channelId string, channelType discordgo.ChannelType, name string, content string) (ch *discordgo.Channel, err error) {
	if channelType == discordgo.ChannelTypeGuildText {
		logging.Info.Printf("Creating a public thread in channel %s\n", channelId)
		ch, err = d.ThreadStartComplex(channelId, &discordgo.ThreadStart{
			Name:                name,
			AutoArchiveDuration: 10080,
			Type:                discordgo.ChannelTypeGuildPublicThread,
		})
		d.ChannelMessageSend(ch.ID, content)

	} else if channelType == discordgo.ChannelTypeGuildForum {
		logging.Info.Printf("Creating a post in channel %s\n", channelId)
		ch, err = d.ForumThreadStartComplex(channelId, &discordgo.ThreadStart{
			Name:                name,
			AutoArchiveDuration: 10080,
		}, &discordgo.MessageSend{
			Content: content,
		})
		logging.Info.Printf("Created a post in channel %s\n", channelId)
	}
	return
}

// Start the discussion about an event with a string message
func (d *discordBot) StartEventDiscussion(ce *models.CosmoEvent, title, message string) (*discordgo.Channel, error) {
	responseMethod, err := ce.GetResponseMethod()

	if err != nil {
		return nil, err
	}

	logging.Debug.Printf("Got response method with channel id%s and method %d\n", responseMethod.DiscordId, responseMethod.ChannelType)

	return d.StartDiscussion(responseMethod.DiscordId, responseMethod.ChannelType, title, message)
}

package models

import (
	"errors"
	"time"

	"github.com/wo0lien/cosmoBot/internal/config"
	"gorm.io/gorm"
)

type CosmoEvent struct {
	gorm.Model
	EventType        config.EventType
	Name             string
	StartDate        time.Time `gorm:"start"`
	EndDate          time.Time `gorm:"end"`
	DoesChannelExist bool      `gorm:"default:false"`
	ChannelID        *string
}

type Volunteer struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Events    []CosmoEvent `gorm:"many2many:volunteers_events;"`
	DiscordID *string
}

type VolunteerEvent struct {
	CosmoEventID           uint `gorm:"primaryKey"`
	VolunteerID            uint `gorm:"primaryKey"`
	VolunteerHasBeenTagged bool `gorm:"default:false"`
}

// Get the channelType of the event discussion on discord
// can be textChannel
// can be forum
func (ce *CosmoEvent) GetResponseMethod() (*config.ResponseMethod, error) {
	if method, ok := config.Config.ResponseMethodByEventType[ce.EventType]; ok {
		return &method, nil
	}
	return nil, errors.New("could not get response method based on the provided type, does it exist in config ?")
}

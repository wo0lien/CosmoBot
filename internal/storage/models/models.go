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

type Volunteers struct {
	gorm.Model
	NocoID int
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

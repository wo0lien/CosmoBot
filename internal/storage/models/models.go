package models

import (
	"errors"
	"time"

	"github.com/wo0lien/cosmoBot/internal/config"
	"gorm.io/gorm"
)

type CosmoEvent struct {
	// gorm.Model is a struct containing fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	gorm.Model

	Name      string
	StartDate time.Time `gorm:"start"`
	EndDate   time.Time `gorm:"end"`

	// volunteers attending the event
	Volunteers []Volunteer `gorm:"many2many:volunteer_events;"`

	// discord related fields

	// Event type used to determine the response method on discord
	EventType config.EventType
	// DoesEventExist is true if the event channel exists on discord
	DoesChannelExist bool `gorm:"default:false"`
	// ChannelID is the id of the event channel on discord
	ChannelID *string

	// google calendar related fields

	// DoesCalendarExist is true if the event exists in the google calendar
	DoesCalendarExist bool `gorm:"default:false"`
	// CalendarID is the id of the event in the google calendar
	CalendarID *string
}

type Volunteer struct {
	// gorm.Model is a struct containing fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	gorm.Model

	FirstName string
	LastName  string
	Email     string
	Phone     string
	// events the volunteer is attending
	Events []CosmoEvent `gorm:"many2many:volunteer_events;"`

	// discord related fields
	DiscordID *string
}

type VolunteerEvent struct {
	CosmoEventID uint `gorm:"primaryKey"`
	VolunteerID  uint `gorm:"primaryKey"`

	// discord related fields

	// HasBeenTagged is true if the volunteer has been tagged on discord
	HasBeenTagged bool `gorm:"default:false"`

	// calendar related fields

	// HasBeenInvited is true if the volunteer has been invited on the calendar event
	HasBeenInvited bool `gorm:"default:false"`
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

// IsUpcoming returns true if the event is upcoming
func (ce *CosmoEvent) IsUpcoming() bool {
	return ce.EndDate.After(time.Now())
}

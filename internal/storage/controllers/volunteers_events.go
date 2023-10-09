package controllers

import (
	"errors"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"gorm.io/gorm"
)

func GetAllVolunteersEvents() *[]models.VolunteerEvent {
	var volunteersEvents []models.VolunteerEvent
	db.DB.Find(&volunteersEvents)

	return &volunteersEvents
}

func SaveVolunteerEvent(volunteerEvent *models.VolunteerEvent) error {
	return db.DB.Save(volunteerEvent).Error
}

func LoadVolunteersEventsJoinFromApi(volunteer *api.VolunteersResponse) (*[]uint, error) {
	var eventIds []uint

	// find volunteer in db
	volunteerInDB := GetVolunteerById(uint(*volunteer.Id))

	// check if volunteer exists in db
	if volunteerInDB == nil {
		return nil, errors.New("volunteer not found in db")
	}

	// Add assocation to the volunteer
	for _, association := range *volunteer.NcCurgNcM2mW5i3lbdpwrs {
		logging.Debug.Printf("Adding association to volunteer %v", association)
		db.DB.Model(volunteerInDB).Association("Events").Append(&models.CosmoEvent{Model: gorm.Model{ID: uint(*association.Table1Id)}})

		eventIds = append(eventIds, uint(*association.Table1Id))
	}

	return &eventIds, nil
}

func LoadVolunteersEventsJoinsFromApi(volunteers *[]api.VolunteersResponse) error {
	for _, volunteer := range *volunteers {
		// get all joins for this volunteer
		joins := GetVolunteersEventsJoinByVolunteerId(uint(*volunteer.Id))

		var eventIds []uint

		for _, join := range *joins {
			eventIds = append(eventIds, join.CosmoEventID)
		}

		ids, err := LoadVolunteersEventsJoinFromApi(&volunteer)
		if err != nil {
			return err
		}

		// check for removed events
		for _, id := range *ids {
			popIdFromList(&eventIds, id)
		}

		for _, id := range eventIds {
			err := DeleteVolunteerEventByVolunteerIdAndEventId(uint(*volunteer.Id), id)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func DeleteVolunteerEventByVolunteerIdAndEventId(volunteerId, eventId uint) error {
	db.DB.Delete(&models.VolunteerEvent{VolunteerID: volunteerId, CosmoEventID: eventId})
	return nil
}

func GetVolunteersEventsJoinByVolunteerId(id uint) *[]models.VolunteerEvent {
	var volunteersEventsJoin []models.VolunteerEvent
	db.DB.Where("volunteer_id = ?", id).Find(&volunteersEventsJoin)
	return &volunteersEventsJoin
}

package controllers

import (
	"time"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"github.com/wo0lien/cosmoBot/internal/utils"
	"gorm.io/gorm"
)

// ########## GET ##########

func AllVolunteersEvents() (*[]models.VolunteerEvent, error) {
	var volunteersEvents []models.VolunteerEvent
	err := db.DB.Find(&volunteersEvents).Error

	return &volunteersEvents, err
}

// get joins for all volunteers for upcoming events
func AllUpcomingVolunteersEvents() (*[]models.VolunteerEvent, error) {
	var volunteersEvents []models.VolunteerEvent
	err := db.DB.
		Joins("JOIN cosmo_events ON volunteer_events.cosmo_event_id = cosmo_events.id").
		Where("cosmo_events.end_date > ?", time.Now()).
		Find(&volunteersEvents).Error

	return &volunteersEvents, err
}

// Get all the events joins for a single volunteer
func AllVolunteerEventsByVolunteerID(volunteerID uint) (*[]models.VolunteerEvent, error) {
	var volunteerEvents []models.VolunteerEvent
	err := db.DB.
		Joins("JOIN volunteers ON volunteer_events.volunteer_id = volunteer.id").
		Where("volunteer.id = ?", volunteerID).
		Find(&volunteerEvents).Error
	return &volunteerEvents, err
}

// Get all the events joins for a single volunteer
func AllVolunteerUpcomingEvents(volunteerID uint) (*[]models.VolunteerEvent, error) {
	var volunteerEvents []models.VolunteerEvent
	err := db.DB.
		Joins("JOIN volunteers ON volunteer_events.volunteer_id = volunteer.id").
		Joins("JOIN cosmo_events ON volunteer_events.cosmo_event_id = cosmo_events.id").
		Where("volunteer.id = ?", volunteerID).
		Find(&volunteerEvents).Error
	return &volunteerEvents, err
}

func VolunteersEventsJoinByVolunteerId(id uint) (*[]models.VolunteerEvent, error) {
	var volunteersEventsJoin []models.VolunteerEvent
	err := db.DB.Where("volunteer_id = ?", id).Find(&volunteersEventsJoin).Error
	return &volunteersEventsJoin, err
}

func VolunteerEventJoinByVolunteerIDAndEventID(volunteerID, eventID uint) (*models.VolunteerEvent, error) {
	var volunteerEvent models.VolunteerEvent
	err := db.DB.Where("volunteer_id = ? AND cosmo_event_id = ?", volunteerID, eventID).First(&volunteerEvent).Error
	return &volunteerEvent, err
}

// ########## UPDATE ##########

// Load every joins a volunteer has in the API
func UpdateVolunteerEventsJoinsFromVolunteerInApi(volunteer *api.VolunteersResponse) (volunteerInDB *models.Volunteer, addedEventsIDs, removedEventsIDs []uint, err error) {
	var eventIDs []uint

	addedEventsIDs = []uint{}
	removedEventsIDs = []uint{}

	// get events from the list of event IDs
	for _, join := range *volunteer.NcCurgNcM2mW5i3lbdpwrs {
		eventIDs = append(eventIDs, uint(*join.Table1Id))
	}

	// get volunteer with events from db
	volunteerInDB, err = VolunteerWithEventsById(uint(*volunteer.Id))
	if err != nil {
		return
	}

	// list events id from db
	var eventIDsInDB []uint
	for _, event := range volunteerInDB.Events {
		eventIDsInDB = append(eventIDsInDB, event.ID)
	}

	// check for new events
	for _, id := range eventIDs {
		if !utils.Contains(eventIDsInDB, id) {
			// check if event exists in db
			eventInDB, err := EventByID(id)
			if err != nil {
				logging.Info.Printf("Event with id %d does not exist in db. Error: %s", id, err)
				continue
			}
			if eventInDB == nil {
				logging.Info.Printf("Event with id %d does not exist in db.", id)
				continue
			}

			// add event to volunteer
			db.DB.Model(volunteerInDB).Association("Events").Append(&models.CosmoEvent{
				Model: gorm.Model{
					ID: id,
				},
			})

			addedEventsIDs = append(addedEventsIDs, id)
		}
	}

	// check for removed events
	for _, id := range eventIDsInDB {
		if !utils.Contains(eventIDs, id) {
			// remove event from volunteer
			db.DB.Model(volunteerInDB).Association("Events").Delete(models.CosmoEvent{
				Model: gorm.Model{
					ID: id,
				},
			})

			removedEventsIDs = append(removedEventsIDs, id)
		}
	}

	return
}

func UpdateVolunteerEventsJoinsFromEventInApi(event *api.EventsResponse) (eventInDB *models.CosmoEvent, addVoluntersIDs, removedVolunteersIDs []uint, err error) {
	var volunteersIDs []uint

	addVoluntersIDs = []uint{}
	removedVolunteersIDs = []uint{}

	// get volunteers from the list of volunteer IDs
	for _, join := range *event.NcCurgNcM2mW5i3lbdpwrs {
		volunteersIDs = append(volunteersIDs, uint(*join.Table2Id))
	}

	eventInDB, err = EventWithVolunteersByID(uint(*event.Id))

	if err != nil {
		return
	}

	// list volunteers id from db
	var volunteersIDsInDB []uint
	for _, volunteer := range eventInDB.Volunteers {
		volunteersIDsInDB = append(volunteersIDsInDB, volunteer.ID)
	}

	// check for new volunteers
	for _, id := range volunteersIDs {
		if !utils.Contains(volunteersIDsInDB, id) {
			// add volunteer to event
			err := db.DB.Model(eventInDB).Association("Volunteers").Append(&models.Volunteer{
				Model: gorm.Model{
					ID: id,
				},
			})

			if err != nil {
				logging.Error.Printf("Could not add volunteer to event. Error: %s", err)
				continue
			}

			addVoluntersIDs = append(addVoluntersIDs, id)
		}
	}

	// check for removed volunteers
	for _, id := range volunteersIDsInDB {
		if !utils.Contains(volunteersIDs, id) {
			// remove volunteer from event
			err := db.DB.Model(eventInDB).Association("Volunteers").Delete(models.Volunteer{
				Model: gorm.Model{
					ID: id,
				},
			})

			if err != nil {
				logging.Error.Printf("Could not delete volunteer from event. Error: %s", err)
				continue
			}

			removedVolunteersIDs = append(removedVolunteersIDs, id)
		}
	}

	return
}

// Wrapper for UpdateVolunteerEventsJoinsFromEventInApi that takes an array of volunteers
func UpdateVolunteersEventsJoinsFromApi(volunteers *[]api.VolunteersResponse) error {
	for _, volunteer := range *volunteers {
		_, _, _, err := UpdateVolunteerEventsJoinsFromVolunteerInApi(&volunteer)
		if err != nil {
			return err
		}
	}
	return nil
}

// ########## SAVE ##########

func SaveVolunteerEvent(volunteerEvent *models.VolunteerEvent) error {
	return db.DB.Save(volunteerEvent).Error
}

// ########## DELETE ##########

func DeleteVolunteerEventByVolunteerIdAndEventId(volunteerId, eventId uint) error {
	db.DB.Delete(&models.VolunteerEvent{VolunteerID: volunteerId, CosmoEventID: eventId})
	return nil
}

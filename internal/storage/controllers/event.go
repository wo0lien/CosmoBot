package controllers

import (
	"errors"
	"time"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/config"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"github.com/wo0lien/cosmoBot/internal/utils"
	"gorm.io/gorm"
)

func SaveEvent(event *models.CosmoEvent) error {
	return db.DB.Save(event).Error
}

func EventByID(id uint) (*models.CosmoEvent, error) {
	var event models.CosmoEvent
	err := db.DB.First(&event, id).Error
	return &event, err
}

func EventWithVolunteersByID(id uint) (*models.CosmoEvent, error) {
	var event models.CosmoEvent
	err := db.DB.Preload("Volunteers").First(&event, id).Error
	return &event, err
}

func AllEvents() (*[]models.CosmoEvent, error) {
	var events []models.CosmoEvent
	err := db.DB.Find(&events).Error
	return &events, err
}

// get all events where end date is greater than today
func AllUpcomingEvents() *[]models.CosmoEvent {
	var events []models.CosmoEvent
	db.DB.Where("end_date > ?", time.Now()).Find(&events)
	return &events
}

// Get all upcoming events ids
func AllUpcomingEventIDs() (*[]uint, error) {
	var eventIDs []uint
	err := db.DB.Select("ID").Where("end_date > ?", time.Now()).Find(&eventIDs).Error
	return &eventIDs, err
}

// Delete CosmoEvent by id in the database
func DeleteEventById(id uint) error {
	db.DB.Delete(&models.CosmoEvent{Model: gorm.Model{ID: id}})
	return nil
}

func DeleteChannelEventByIdAndSave(id uint) error {
	var event models.CosmoEvent
	err := db.DB.Model(event).Where("channel_id = ?", id).Update("channel_id", nil).Update("does_channel_exist", false).Error
	if err != nil {
		return err
	}
	err = SaveEvent(&event)

	if err != nil {
		return err
	}

	return nil
}

// load an event in the database from the api response format
func CreateOrUpdateEventInDBFromApi(event api.EventsResponse) (*models.CosmoEvent, error) {

	logging.Info.Println("Got an event to load into DB")

	if event.Type == nil {
		return nil, errors.New("event type is nil, could not load type in db")
	}

	StartDate, err := time.Parse(api.NOCO_TIME_LAYOUT, *event.Start)
	if err != nil {
		return nil, err
	}

	EndDate, err := time.Parse(api.NOCO_TIME_LAYOUT, *event.End)
	if err != nil {
		return nil, err
	}

	eventInDB := models.CosmoEvent{
		EventType: config.EventType(*event.Type),
		Name:      *event.Title, // never nil
		StartDate: StartDate,
		EndDate:   EndDate,
		Model: gorm.Model{
			ID: uint(*event.Id), // never nil
		},
	}

	err = db.DB.Save(&eventInDB).Error
	if err != nil {
		return nil, err
	}

	logging.Info.Println("Loaded event into DB")

	return &eventInDB, nil
}

// load a list of events in the database from the api response format
func CreateOrUpdateEventsInDBFromApi(events []api.EventsResponse) error {
	// store a list of Ids
	var ids []uint
	// load events ids in the list
	eventsFromDB, err := AllEvents()
	if err != nil {
		return err
	}
	for _, event := range *eventsFromDB {
		ids = append(ids, event.ID)
	}

	for _, event := range events {
		logging.Debug.Printf("Loading event with id %d\n in db", *event.Id)
		_, err := CreateOrUpdateEventInDBFromApi(event)
		if err != nil {
			logging.Warning.Printf("error loading the event %s : %s ", *event.Title, err)
			continue
		}
		// remove the id from the list
		utils.PopIdFromList(&ids, uint(*event.Id))
	}

	for _, el := range ids {
		logging.Debug.Printf("Deleting event with id %d\n", el)
		DeleteEventById(uint(el))
	}

	return nil
}

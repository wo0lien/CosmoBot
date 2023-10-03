package controllers

import (
	"errors"
	"time"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/config"
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"gorm.io/gorm"
)

func SaveEvent(event *models.CosmoEvent) {
	db.DB.Save(event)
}

func GetEventByID(id uint) *models.CosmoEvent {
	var event models.CosmoEvent
	db.DB.First(&event, id)
	return &event
}

func GetAllCosmoEvents() *[]models.CosmoEvent {
	var events []models.CosmoEvent
	db.DB.Find(&events)
	return &events
}

// get all events where end date is greater than today
func GetAllUpcomingCosmoEvents() *[]models.CosmoEvent {
	var events []models.CosmoEvent
	db.DB.Where("end_date > ?", time.Now()).Find(&events)
	return &events
}

// Delete CosmoEvent by id in the database
func DeleteEventById(id uint) error {
	db.DB.Delete(&models.CosmoEvent{Model: gorm.Model{ID: id}})
	return nil
}

// load an event in the database from the api response format
func LoadEventInDBFromAPI(event api.EventsResponse) error {

	logging.Info.Println("Got an event to load into DB")

	if event.Type == nil {
		return errors.New("event type is nil, could not load type in db")
	}
	StartDate, err := time.Parse(api.NOCO_TIME_LAYOUT, *event.Start)
	if err != nil {
		return err
	}
	EndDate, err := time.Parse(api.NOCO_TIME_LAYOUT, *event.End)
	if err != nil {
		return err
	}

	db.DB.Create(&models.CosmoEvent{
		EventType: config.EventType(*event.Type),
		Name:      *event.Title, // never nil
		StartDate: StartDate,
		EndDate:   EndDate,
		Model: gorm.Model{
			ID: uint(*event.Id),
		},
	})

	logging.Info.Println("Loaded event into DB")

	return nil
}

// remove an int from a list of ints
// change the order of the element for performance reasons
func popIdFromList(ids *[]uint, id uint) *[]uint {
	for i, el := range *ids {
		if el == id {
			// remove the id from the list
			// order does not matter
			(*ids)[i] = (*ids)[len(*ids)-1]
			*ids = (*ids)[:len(*ids)-1]
			break
		}
	}
	return ids
}

// load a list of events in the database from the api response format
func LoadEventsInDBFromAPI(events []api.EventsResponse) error {
	// store a list of Ids
	var ids []uint
	// load events ids in the list
	for _, event := range *GetAllCosmoEvents() {
		ids = append(ids, event.ID)
	}

	for _, event := range events {
		logging.Debug.Printf("Loading event with id %d\n in db", *event.Id)
		err := LoadEventInDBFromAPI(event)
		if err != nil {
			return err
		}
		// remove the id from the list
		popIdFromList(&ids, uint(*event.Id))
	}

	for _, el := range ids {
		logging.Debug.Printf("Deleting event with id %d\n", el)
		DeleteEventById(uint(el))
	}

	return nil
}

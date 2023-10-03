package controllers

import (
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
)

func GetAllVolunteersEvents() *[]models.VolunteerEvent {
	var volunteersEvents []models.VolunteerEvent
	db.DB.Find(&volunteersEvents)

	return &volunteersEvents
}

func SaveVolunteerEvent(volunteerEvent *models.VolunteerEvent) {
	db.DB.Save(volunteerEvent)
}

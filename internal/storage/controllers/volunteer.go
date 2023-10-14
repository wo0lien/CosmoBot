package controllers

import (
	"errors"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/storage/db"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"gorm.io/gorm"
)

func GetAllVolunteers() *[]models.Volunteer {
	var volunteers []models.Volunteer

	db.DB.Find(&volunteers)

	return &volunteers
}

func SaveVolunteer(volunteer *models.Volunteer) {
	db.DB.Save(volunteer)
}

func DeleteVolunteerById(id uint) error {
	db.DB.Delete(&models.Volunteer{Model: gorm.Model{ID: id}})
	return nil
}

func VolunteerById(id uint) (*models.Volunteer, error) {
	var volunteer models.Volunteer
	err := db.DB.Preload("Events").First(&volunteer, id).Error
	return &volunteer, err
}

func VolunteerWithEventsById(id uint) (*models.Volunteer, error) {
	var volunteer models.Volunteer
	err := db.DB.Preload("Events").First(&volunteer, id).Error
	return &volunteer, err
}

func VolunteerByDiscordId(id string) *models.Volunteer {
	var volunteer models.Volunteer
	db.DB.Where("discord_id = ?", id).First(&volunteer)
	return &volunteer
}

func AllVolunteersByEventID(id uint) (*[]models.Volunteer, error) {
	var volunteers []models.Volunteer
	err := db.DB.Model(&models.CosmoEvent{Model: gorm.Model{ID: id}}).Association("Volunteers").Find(&volunteers)
	return &volunteers, err
}

func CreateOrUpdateVolunteerToDBFromAPI(volunteer *api.VolunteersResponse) (*models.Volunteer, error) {
	// check if pointers to field are not empty
	if volunteer.Firstname == nil || volunteer.Lastname == nil || volunteer.Email == nil || volunteer.Tel == nil {
		return nil, errors.New("one of the fields is nil, could not load volunteer in db")
	}

	vol := models.Volunteer{
		Model:     gorm.Model{ID: uint(*volunteer.Id)}, // never nil
		FirstName: *volunteer.Firstname,
		LastName:  *volunteer.Lastname,
		Email:     *volunteer.Email,
		Phone:     *volunteer.Tel,
		DiscordID: volunteer.DiscordId,
	}

	err := db.DB.Save(&vol).Error

	return &vol, err
}

func CreateOrUpdateVolunteersToDBFromAPI(volunteers *[]api.VolunteersResponse) error {
	for _, volunteer := range *volunteers {
		_, err := CreateOrUpdateVolunteerToDBFromAPI(&volunteer)
		if err != nil {
			return err
		}
	}
	return nil
}

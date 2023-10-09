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

func GetVolunteerById(id uint) *models.Volunteer {
	var volunteer models.Volunteer
	db.DB.Preload("Events").First(&volunteer, id)
	return &volunteer
}

func GetVolunteerByDiscordId(id string) *models.Volunteer {
	var volunteer models.Volunteer
	db.DB.Where("discord_id = ?", id).First(&volunteer)
	return &volunteer
}

func LoadVolunteerToDBFromAPI(volunteer *api.VolunteersResponse) error {
	// check if pointers to field are not empty
	if volunteer.Firstname == nil || volunteer.Lastname == nil || volunteer.Email == nil || volunteer.Tel == nil {
		return errors.New("one of the fields is nil, could not load volunteer in db")
	}

	db.DB.Create(&models.Volunteer{
		Model:     gorm.Model{ID: uint(*volunteer.Id)}, // never nil
		FirstName: *volunteer.Firstname,
		LastName:  *volunteer.Lastname,
		Email:     *volunteer.Email,
		Phone:     *volunteer.Tel,
		DiscordID: volunteer.DiscordId,
	})

	return nil
}

func LoadVolunteersToDBFromAPI(volunteers *[]api.VolunteersResponse) error {
	for _, volunteer := range *volunteers {
		err := LoadVolunteerToDBFromAPI(&volunteer)
		if err != nil {
			return err
		}
	}
	return nil
}

package controllers

import (
	"errors"

	"github.com/wo0lien/cosmoBot/internal/api"
	"github.com/wo0lien/cosmoBot/internal/logging"
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

func LoadVolunteersEventsJoinsFromApi(volunteers *[]api.VolunteersResponse) error {
	for _, volunteer := range *volunteers {
		err := LoadVolunteerEventsJoinsFromApi(&volunteer)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadVolunteerEventsJoinsFromApi(volunteer *api.VolunteersResponse) error {
	// find volunteer in db
	volunteerInDB := GetVolunteerById(uint(*volunteer.Id))

	// check if volunteer exists in db
	if volunteerInDB == nil {
		return errors.New("volunteer not found in db")
	}

	// Add assocation to the volunteer
	for _, association := range *volunteer.NcCurgNcM2mW5i3lbdpwrs {
		logging.Debug.Printf("Adding association to volunteer %v", association)
		db.DB.Model(volunteerInDB).Association("Events").Append(&models.CosmoEvent{Model: gorm.Model{ID: uint(*association.Table1Id)}})
	}

	return nil
}

func GetVolunteersEventsJoinByVolunteerId(id uint) *[]models.VolunteerEvent {
	var volunteersEventsJoin []models.VolunteerEvent
	db.DB.Where("volunteer_id = ?", id).Find(&volunteersEventsJoin)
	return &volunteersEventsJoin
}

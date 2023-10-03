package db

import (
	"github.com/wo0lien/cosmoBot/internal/logging"
	"github.com/wo0lien/cosmoBot/internal/storage/models"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

type DbStruct struct {
	*gorm.DB
}

var DB DbStruct

func init() {
	logging.Info.Println("Connecting to database")
	db, err := Connect()

	if err != nil {
		logging.Critical.Fatalf("Could not connect to database, error : %s", err)
	}

	DB = DbStruct{db}

	logging.Info.Println("Migrating the models")
	DB.AutoMigrate(&models.CosmoEvent{})

	logging.Info.Println("Database ready")
}

// Connect to the database
func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

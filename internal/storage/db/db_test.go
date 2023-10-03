package db_test

import (
	"testing"

	"github.com/wo0lien/cosmoBot/internal/storage/db"
)

func TestConnect(t *testing.T) {
	db, err := db.Connect()

	if err != nil {
		t.Error(err)
	}

	if db == nil {
		t.Error("db is nil")
	}
}

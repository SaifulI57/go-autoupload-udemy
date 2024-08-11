package store_test

import (
	"testing"
	"time"

	"github.com/SaifulI57/uploader-udemy/store"
)

func TestGetAllDB(t *testing.T) {
	_, err := upload.GetAllDBUploader(store.Store)
	if err != nil {
		t.Error("Failed to get uploader")
	}
}

func TestGetAllDateDB(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-08-08")
	_, err := upload.GetAllSDateDBUploader(store.Store, date)
	if err != nil {
		t.Error("Error retrieving uploader with specified date")
	}
}

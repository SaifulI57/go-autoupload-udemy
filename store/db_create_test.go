package store_test

import (
	"fmt"
	"testing"
	"time"

	log "github.com/SaifulI57/uploader-udemy/logger"
	"github.com/SaifulI57/uploader-udemy/store"
)

var upload *store.UploadDB

func init() {
	store.ConnectDB()
	upload = &store.UploadDB{}
}

func TestCreateDB(t *testing.T) {

	newUploader := &store.UploaderData{
		UID:                1,
		UdemyUrl:           "stesting.com7",
		Description:        "testing7",
		HashUrl:            "tesing7",
		UploadedToFacebook: false,
		UploadedToBot:      false,
		UploadedToWebsite:  false,
		UploadedToWhatsapp: false,
		ExpiredIn:          time.Now(),
	}
	if upload.CreateDBUploader(store.Store, newUploader) != nil {
		t.Error("Failed to create uploader")
	}

}

func TestCreateBatchDB(t *testing.T) {
	newUploader := []*store.UploaderData{{
		UID:                2,
		UdemyUrl:           "stesting.com1",
		Description:        "testing",
		HashUrl:            "tesing1",
		UploadedToFacebook: false,
		UploadedToBot:      false,
		UploadedToWebsite:  false,
		UploadedToWhatsapp: false,
		ExpiredIn:          time.Now(),
	}, {
		UID:                3,
		UdemyUrl:           "stesting.com2",
		Description:        "testing",
		HashUrl:            "tesing2",
		UploadedToFacebook: false,
		UploadedToBot:      false,
		UploadedToWebsite:  false,
		UploadedToWhatsapp: false,
		ExpiredIn:          time.Now(),
	}, {
		UID:                4,
		UdemyUrl:           "stesting.com3",
		Description:        "testing",
		HashUrl:            "tesing3",
		UploadedToFacebook: false,
		UploadedToBot:      false,
		UploadedToWebsite:  false,
		UploadedToWhatsapp: false,
		ExpiredIn:          time.Now(),
	}, {
		UID:                5,
		UdemyUrl:           "stesting.com4",
		Description:        "testing",
		HashUrl:            "tesing4",
		UploadedToFacebook: false,
		UploadedToBot:      false,
		UploadedToWebsite:  false,
		UploadedToWhatsapp: false,
		ExpiredIn:          time.Now(),
	}}
	for _, upload := range newUploader {
		log.Logger.Info(fmt.Sprintf("Uploading UID: %d", upload.UID))
	}
	if upload.BatchUploadDBUploader(store.Store, newUploader) != nil {
		t.Error("Failed to Create Batch")
	}
}

package store_test

import (
	"testing"

	"github.com/SaifulI57/uploader-udemy/store"
)

func TestDeleteDB(t *testing.T) {
	newUploader := &store.UploaderData{
		UID: 1,
	}
	if upload.DeleteDBUploader(store.Store, newUploader) != nil {
		t.Error("Failed to delete uploader")
	}
}

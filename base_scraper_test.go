package main_test

import (
	"log"
	"testing"

	m "github.com/SaifulI57/uploader-udemy"
	"github.com/stretchr/testify/assert"
)

func TestLiveScraper(t *testing.T) {
	d := &m.Doscrape{}
	result, err := d.ListToday()

	if err != nil {
		t.Fatalf("Scraping failed: %v", err)
	}

	assert.NotNil(t, result)
	assert.True(t, len(result) > 0, "The scraper should return some data")

	// Log the first few results for verification
	for i, res := range result {
		if i >= 5 { // limit output for test readability
			break
		}
		log.Printf("Result %d: %v\n", i+1, res)
	}
}

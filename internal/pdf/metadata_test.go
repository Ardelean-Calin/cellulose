package pdf

import (
	"os"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func TestGetCreationDate(t *testing.T) {
	// Create a temporary PDF file
	tmpFile := "test_creation_date.pdf"
	defer os.Remove(tmpFile)

	// Set an arbitrary creation date
	creationDate := time.Date(2025, time.February, 2, 15, 4, 5, 0, time.UTC)

	// Create an empty PDF with the specified creation date
	conf := pdfcpu.NewDefaultConfiguration()
	conf.TimeNow = func() time.Time { return creationDate }

	err := api.CreatePDFFile([]string{}, nil, tmpFile, conf)
	if err != nil {
		t.Fatalf("Failed to create PDF file: %v", err)
	}

	// Use the function to read back the creation date
	extractedDateStr, err := GetCreationDate(tmpFile)
	if err != nil {
		t.Fatalf("Failed to get creation date: %v", err)
	}

	// Parse the extracted date
	extractedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", extractedDateStr)
	if err != nil {
		t.Fatalf("Failed to parse creation date: %v", err)
	}

	// Compare the dates
	if !creationDate.Equal(extractedDate) {
		t.Errorf("Expected creation date %v, got %v", creationDate, extractedDate)
	}
}

package pdf

import (
	"os"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestGetCreationDate(t *testing.T) {
	// Create a temporary PDF file
	tmpFile := "test_creation_date.pdf"
	defer os.Remove(tmpFile)

	// Set an arbitrary creation date
	creationDate := time.Date(2025, time.February, 2, 15, 4, 5, 0, time.UTC)

	table, err := pdfcpu.CreateDemoXRef()
	if err != nil {
		t.Fatalf("Failed to create PDF file table: %v", err)
	}

	err = api.CreatePDFFile(table, tmpFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("Failed to create PDF file: %v", err)
	}

	// Use the function to read back the creation date
	extractedDate, err := GetCreationDate(tmpFile)
	if err != nil {
		t.Fatalf("Failed to get creation date: %v", err)
	}

	// Compare the dates
	if !creationDate.Equal(extractedDate) {
		t.Errorf("Expected creation date %v, got %v", creationDate, extractedDate)
	}
}

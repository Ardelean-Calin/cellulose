package pdf

import (
	"testing"
	"time"
)

func TestGetCreationDate(t *testing.T) {
	expected := []struct {
		pdfFile string
		date    time.Time
	}{
		{
			pdfFile: "testdata/test1.pdf",
			date:    time.Date(2024, 10, 3, 17, 23, 19, 0, time.UTC),
		},
		{
			pdfFile: "testdata/test2.pdf",
			date:    time.Date(2025, 1, 31, 11, 07, 27, 0, time.UTC),
		},
	}

	for _, e := range expected {
		extractedDate, err := GetCreationDate(e.pdfFile)
		if err != nil {
			t.Errorf("Failed to get creation date from %s: %v", e.pdfFile, err)
		}

		if !extractedDate.Equal(e.date) {
			t.Errorf("Dates don't match: Expected: %s, Got: %s", e.date, extractedDate)
		}

	}
}

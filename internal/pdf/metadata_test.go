package pdf

import (
	"os"
	"testing"
	"time"

)

func TestGetCreationDate(t *testing.T) {



	testFiles := []string{"testdata/test1.pdf", "testdata/test2.pdf"}

	for _, file := range testFiles {
		extractedDate, err := GetCreationDate(file)
		if err != nil {
			t.Errorf("Failed to get creation date from %s: %v", file, err)
			continue
		}

		if extractedDate.IsZero() {
			t.Errorf("Expected a valid creation date for %s, got zero value", file)
		}
	}
}

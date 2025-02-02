package pdf

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

// GetCreationDate extracts the creation date from a PDF's metadata
func GetCreationDate(filePath string) (time.Time, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Updated regex to handle both 'Z' and timezone offset formats
	// Now captures everything between the tags that looks like a timestamp
	dateRegex := regexp.MustCompile(`<xmp:CreateDate>([^<]+)</xmp:CreateDate>`)

	// Scan through the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line matches our pattern
		matches := dateRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			// Try parsing with different time formats
			dateStr := matches[1]

			// First try RFC3339 format (handles both Z and timezone offset)
			t, err := time.Parse(time.RFC3339, dateStr)
			if err == nil {
				return t.UTC(), nil
			}

			// If that fails, try additional formats if needed
			// Add more formats here if you encounter other variations

			// If all parsing attempts fail
			return time.Time{}, fmt.Errorf("unable to parse date '%s': %v", dateStr, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("error reading file: %v", err)
	}

	return time.Time{}, fmt.Errorf("date not found in file")
}

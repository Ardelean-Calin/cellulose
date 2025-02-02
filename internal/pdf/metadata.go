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
	// Updated regex to handle both XML and PDF trailer creation dates
	dateRegex := regexp.MustCompile(`<xmp:CreateDate>([^<]+)</xmp:CreateDate>|/CreationDate\s*\(([^)]+)\)`)

	// Scan through the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line matches our pattern
		matches := dateRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			// Determine which capture group has the date string
			var dateStr string
			if matches[1] != "" {
				dateStr = matches[1]
				// Try parsing date in RFC3339 format
				t, err := time.Parse(time.RFC3339, dateStr)
				if err == nil {
					return t.UTC(), nil
				}
			} else if matches[2] != "" {
				dateStr = matches[2]
				// Try parsing date in various PDF date formats
				layouts := []string{
					"D:20060102150405",           // Without timezone
					"D:20060102150405Z07'00'",    // With timezone offset
					"D:20060102150405-07'00'",    // With negative timezone offset
					"D:20060102150405Z",          // With 'Z' timezone
				}
				for _, layout := range layouts {
					t, err := time.Parse(layout, dateStr)
					if err == nil {
						return t.UTC(), nil
					}
				}
				// Try parsing date in 'M/D/YYYY HH:MM:SS' format
				t, err := time.Parse("1/2/2006 15:04:05", dateStr)
				if err == nil {
					return t.UTC(), nil
				}
			}

			// If parsing fails, continue to next line
		}
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("error reading file: %v", err)
	}

	return time.Time{}, fmt.Errorf("date not found in file")
}

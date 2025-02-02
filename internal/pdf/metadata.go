package pdf

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

// GetCreationDate extracts the creation date from a PDF's metadata
func GetCreationDate(filePath string) (string, error) {
	ctx, err := api.ReadContextFile(filePath, pdfcpu.NewDefaultConfiguration())
	if err != nil {
		return "", fmt.Errorf("failed to read PDF context: %w", err)
	}

	infoDict, err := ctx.Info()
	if err != nil {
		return "", fmt.Errorf("failed to get info dictionary: %w", err)
	}

	creationDate, found := infoDict.DateEntry("CreationDate")
	if !found {
		return "", fmt.Errorf("creation date not found in metadata")
	}

	return creationDate.Format("2006-01-02T15:04:05Z07:00"), nil
}

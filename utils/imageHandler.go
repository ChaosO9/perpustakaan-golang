package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ImageUpload uploads an image file.
func ImageUpload(c *gin.Context, formFileKey, uploadDir string) (string, error) {
	file, err := c.FormFile(formFileKey)
	if err != nil {
		return "", fmt.Errorf("failed to get form file: %w", err)
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		return "", fmt.Errorf("failed to save uploaded file: %w", err)
	}

	return filename, nil
}


// ImageDelete deletes an image file.
func ImageDelete(uploadDir, filename string) error {
	filepath := filepath.Join(uploadDir, filename)
	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}
	return nil
}



// ImageUpdate updates an image file. It deletes the old image if a new one is provided.
func ImageUpdate(c *gin.Context, formFileKey, uploadDir, oldFilename string) (string, error) {
    _, err := c.FormFile(formFileKey)
    if err != nil && err != http.ErrMissingFile { // Error other than missing file
        return "", fmt.Errorf("failed to get form file: %w", err) 
    }

    if err == nil { // New file provided
        // Delete old image if it exists
        if oldFilename != "" {
            if err := ImageDelete(uploadDir, oldFilename); err != nil {
                return "", fmt.Errorf("failed to delete old image: %w", err)
            }
        }

        // Upload new image
        newFilename, uploadErr := ImageUpload(c, formFileKey, uploadDir)
        if uploadErr != nil {
            return "", uploadErr
        }
        return newFilename, nil
    } else if err == http.ErrMissingFile { // No new file, keep old filename
        return oldFilename, nil 
    }
    return "", err // Shouldn't reach here, but handle for completeness
}
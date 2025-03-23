package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
)

// ✅ Generate a Secure Random Filename
func GenerateSecureFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "default_filename" + ext
	}
	return fmt.Sprintf("%s%s", hex.EncodeToString(randomBytes), ext)
}

package helpers

import (
	"fmt"
	"github.com/koshka_backend/middleware"
	"math/rand"
	"time"
)

// GenerateUniqueCode creates a random 5-digit code ensuring uniqueness
func GenerateUniqueCode() (string, error) {
	rand.Seed(time.Now().UnixNano())

	// Step 1: Get all active codes in cache
	existingCodes := middleware.CacheListKeys("code_") // Fetch all keys with prefix "code_"
	usedCodes := make(map[string]bool)

	for _, key := range existingCodes {
		usedCodes[key[5:]] = true // Remove "code_" prefix
	}

	// Step 2: Find a free 5-digit code
	if len(usedCodes) >= 100000 {
		return "", fmt.Errorf("all possible 5-digit codes are in use")
	}

	for {
		code := fmt.Sprintf("%05d", rand.Intn(100000)) // Generate 5-digit code
		if !usedCodes[code] {                          // ✅ If not in map, it's unique
			return code, nil
		}
	}
}

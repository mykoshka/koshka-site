package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/koshka_backend/config"
	"github.com/koshka_backend/helpers"
	"github.com/koshka_backend/middleware"
	"log"
	"time"
)

// CollarResponse defines the response structure excluding the phone number
type CollarResponse struct {
	TagID            string `json:"tag_id"`
	PetName          string `json:"pet_name"`
	PetImage         string `json:"pet_image"`
	PetDOB           string `json:"pet_dob"`
	PetNeutered      bool   `json:"neutered"`
	PetVaccinated    bool   `json:"vaccinated"`
	VerificationCode string `json:"extention"`
}

// Response structure for valid verification
type VerifyCodeResponse struct {
	PhoneNumber string `json:"phone_number"`
}

// FetchCollarInfo retrieves collar, pet details, and caches a 5-digit verification code
func FetchCollarInfo(c *fiber.Ctx) error {
	tagID := c.Params("tag_id")

	// Fetch Collar Info & Pet Details
	var collarID, petName, petImage, petDOB string
	var petNeutered, petVaccinated bool
	var phoneNumber string

	err := config.DB.QueryRow(`
    SELECT rc.tag_id, pp.name, img.file, pp.date_of_birth, pp.neutered, pp.vaccinated, u.mobile_number
    FROM reunite_collars rc
    JOIN pet_collars pc ON rc.tag_id = pc.collar_id  -- ✅ Link collar to pet
    JOIN pet_profiles pp ON pc.pet_id = pp.id  -- ✅ Get pet details
    JOIN user_pets up ON pp.id = up.pet_id  -- ✅ Link pet to user
    JOIN user_profiles u ON up.user_id = u.email  -- ✅ Get user details
    LEFT JOIN images img ON pp.picture_id = img.id
    WHERE rc.tag_id = $1 AND rc.registered = TRUE`, tagID).
		Scan(&collarID, &petName, &petImage, &petDOB, &petNeutered, &petVaccinated, &phoneNumber)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collar not found or not registered"})
	} else if err != nil {
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while fetching collar info"})
	}

	// Generate a Temporary 5-Digit Code
	verificationCode, err := helpers.GenerateUniqueCode()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Store code and phone number in cache (expires in 15 minutes)
	cacheKey := "code_" + verificationCode
	cacheValue, _ := json.Marshal(fiber.Map{"phone_number": phoneNumber})
	middleware.CacheSet(cacheKey, string(cacheValue), 15*time.Minute)

	// Return Collar & Pet Info (Exclude Phone Number)
	return c.Status(fiber.StatusOK).JSON(CollarResponse{
		TagID:            collarID,
		PetName:          petName,
		PetImage:         petImage,
		PetDOB:           petDOB,
		PetNeutered:      petNeutered,
		PetVaccinated:    petVaccinated,
		VerificationCode: verificationCode,
	})
}

// ✅ Fetch phone number by 5-digit code
func FetchPhoneNumberByCode(c *fiber.Ctx) error {

	// ✅ Parse incoming request
	code := c.Params("code")
	cacheKey := "code_" + code

	// ✅ Check if the code exists in cache
	cacheData, exists := middleware.CacheGet(cacheKey)
	if !exists {
		fmt.Println("❌ FetchPhoneNumberByCode: Code not found in cache ->", code)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid or expired code"})
	}

	// ✅ Parse the cached JSON data
	var phoneData VerifyCodeResponse
	if err := json.Unmarshal([]byte(cacheData), &phoneData); err != nil {
		fmt.Println("❌ FetchPhoneNumberByCode: Failed to parse cache data ->", cacheData)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve phone number"})
	}

	// ✅ Return only the phone number
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"phone_number": phoneData.PhoneNumber})
}

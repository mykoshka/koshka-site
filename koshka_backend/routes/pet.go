package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/koshka_backend/config"
	"github.com/koshka_backend/helpers"
	"log"
	"os"
	"strings"
)

// AssignCollarRequest represents the request payload for assigning a collar
type AssignCollarRequest struct {
	TagID string `json:"tag_id"`
}

// AddPetRequest Struct for pet request payload
type AddPetRequest struct {
	Name        string `json:"name"`
	BirthDate   string `json:"birth_date"`
	TagID       string `json:"tag_id"`
	CityLicence string `json:"city_licence"`
	Neutered    bool   `json:"neutered"`
	Vaccinated  bool   `json:"vaccinated"`
}

// AddPet Add Pet with Transaction Support
func AddPet(c *fiber.Ctx) error {
	userEmail, ok := c.Locals("user").(string)
	if !ok || userEmail == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
	}

	file := form.File["image"][0]
	newFile := helpers.GenerateSecureFilename(file.Filename)
	tempPath := "./uploads/temp_" + file.Filename

	// ✅ Save image temporarily
	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
	}

	// ✅ Generate SHA256 hash
	hash, err := helpers.GenerateImageHash(tempPath)
	if err != nil {
		err := os.Remove(tempPath)
		if err != nil {
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to compute image hash"})
	}

	// ✅ Default `date_of_birth`
	birthDate := c.FormValue("birth_date")
	if birthDate == "" {
		birthDate = "1900-01-01"
	}

	// ✅ Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction error"})
	}

	// Rollback if not committed
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	// ✅ Check if image exists
	var imageID int
	var existingImageFile string
	var newFilePath string

	err = tx.QueryRow("SELECT id, file FROM images WHERE hash = $1 LIMIT 1", hash).Scan(&imageID, &existingImageFile)

	if err == nil {
		// ✅ Duplicate found, use existing image
		err := os.Remove(tempPath)
		if err != nil {
		}
		fmt.Println("⚠️ Duplicate image detected ->", existingImageFile)
	} else if errors.Is(err, sql.ErrNoRows) {
		// ✅ Move new image to permanent storage
		newFilePath = "./uploads/" + newFile
		if err := os.Rename(tempPath, newFilePath); err != nil {
			log.Println("Error moving image to permanent storage", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store image"})
		}

		// ✅ Insert new image
		err = tx.QueryRow("INSERT INTO images (file, hash) VALUES ($1, $2) RETURNING id", newFilePath, hash).Scan(&imageID)
		if err != nil {
			err := os.Remove(tempPath)
			if err != nil {
			}
			log.Println("Error saving image metadata:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image metadata"})
		}
	} else {
		// ✅ Handle unexpected DB errors
		fmt.Println("❌ Database error:", err)
		err := os.Remove(tempPath)
		if err != nil {
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Insert Pet
	var petID string
	err = tx.QueryRow(`
		INSERT INTO pet_profiles (name, date_of_birth, city_licence, neutered, vaccinated, picture_id) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		c.FormValue("name"), birthDate, c.FormValue("city_licence"),
		c.FormValue("neutered") == "true", c.FormValue("vaccinated") == "true",
		imageID,
	).Scan(&petID)

	if err != nil {
		log.Println("Error failed to create pet profile: ", err)
		err := os.Remove(newFilePath)
		if err != nil {
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create pet profile"})
	}

	// ✅ Insert Tag into user_pets
	_, err = tx.Exec("INSERT INTO user_pets (user_id, pet_id) VALUES ($1, $2)",
		userEmail, petID)
	if err != nil {
		log.Println("Error failed to link tag: ", err)
		err := os.Remove(newFilePath)
		if err != nil {
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to link tag to pet"})
	}

	// ✅ Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("Error ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pet added successfully", "pet_id": petID})
}

// UpdatePet - Only Updates Supplied Fields
func UpdatePet(c *fiber.Ctx) error {
	userEmail, ok := c.Locals("user").(string)
	if !ok || userEmail == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	petID := c.Params("id")

	// ✅ Start Transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction error"})
	}

	// Rollback if not committed
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	// ✅ Prepare Query for Updates (Only Provided Fields)
	var updateFields []string
	var values []interface{}
	paramIndex := 1

	// ✅ Check Each Possible Field
	if name := c.FormValue("name"); name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", paramIndex))
		values = append(values, name)
		paramIndex++
	}
	if birthDate := c.FormValue("birth_date"); birthDate != "" {
		updateFields = append(updateFields, fmt.Sprintf("birth_date = $%d", paramIndex))
		values = append(values, birthDate)
		paramIndex++
	}
	if tagID := c.FormValue("tag_id"); tagID != "" {
		updateFields = append(updateFields, fmt.Sprintf("tag_id = $%d", paramIndex))
		values = append(values, tagID)
		paramIndex++
	}
	if cityLicence := c.FormValue("city_licence"); cityLicence != "" {
		updateFields = append(updateFields, fmt.Sprintf("city_licence = $%d", paramIndex))
		values = append(values, cityLicence)
		paramIndex++
	}
	if neutered := c.FormValue("neutered"); neutered != "" {
		updateFields = append(updateFields, fmt.Sprintf("neutered = $%d", paramIndex))
		values = append(values, neutered == "true")
		paramIndex++
	}
	if vaccinated := c.FormValue("vaccinated"); vaccinated != "" {
		updateFields = append(updateFields, fmt.Sprintf("vaccinated = $%d", paramIndex))
		values = append(values, vaccinated == "true")
		paramIndex++
	}

	// ✅ Handle Image Upload (Check for Duplicates)
	form, err := c.MultipartForm()
	if err == nil && len(form.File["image"]) > 0 {
		file := form.File["image"][0]
		tempPath := "./uploads/temp_" + file.Filename

		// ✅ Save Temporarily
		if err := c.SaveFile(file, tempPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image"})
		}

		// ✅ Generate SHA256 Hash
		hash, err := helpers.GenerateImageHash(tempPath)
		if err != nil {
			err := os.Remove(tempPath)
			if err != nil {
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to compute image hash"})
		}

		// ✅ Check if Image Exists in DB
		var existingImageID int
		var existingImageFile string
		err = tx.QueryRow("SELECT id, file FROM images WHERE hash = $1 LIMIT 1", hash).Scan(&existingImageID, &existingImageFile)

		if err == nil {
			// ✅ Duplicate Found, Use Existing Image
			err := os.Remove(tempPath)
			if err != nil {
			}
			fmt.Println("⚠️ Duplicate image detected ->", existingImageFile)
		} else if errors.Is(err, sql.ErrNoRows) {
			// ✅ Move New Image to Permanent Storage
			newFilePath := "./uploads/" + file.Filename
			if err := os.Rename(tempPath, newFilePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to store image"})
			}

			// ✅ Insert New Image into DB
			err = tx.QueryRow("INSERT INTO images (file, hash) VALUES ($1, $2) RETURNING id", newFilePath, hash).Scan(&existingImageID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save image metadata"})
			}
		} else {
			// ✅ Handle Unexpected DB Errors
			fmt.Println("❌ Database error:", err)
			err := os.Remove(tempPath)
			if err != nil {
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		// ✅ Add Image Update to Query
		updateFields = append(updateFields, fmt.Sprintf("image_id = $%d", paramIndex))
		values = append(values, existingImageID)
		paramIndex++
	}

	// ✅ Handle `tag_id` Separately in `user_pets`
	if tagID := c.FormValue("tag_id"); tagID != "" {
		var existingTagID string
		err := tx.QueryRow("SELECT tag_id FROM reunite_collars WHERE tag_id = $1", tagID).Scan(&existingTagID)

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		_, err = tx.Exec("UPDATE user_pets SET tag_id = $1 WHERE pet_id = $2", tagID, petID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update tag ID"})
		}
	}

	// ✅ Finalize Query
	if len(updateFields) > 0 {
		query := fmt.Sprintf("UPDATE pet_profiles SET %s, updated_at = NOW() WHERE id = $%d",
			strings.Join(updateFields, ", "), paramIndex)
		values = append(values, petID)

		stmt, err := tx.Prepare(query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)

		_, err = stmt.Exec(values...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update pet"})
		}
	}

	// ✅ Commit Transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pet updated successfully"})
}

// DeletePet - Securely deletes a pet
func DeletePet(c *fiber.Ctx) error {
	userEmail, ok := c.Locals("user").(string)
	if !ok || userEmail == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse pet ID from URL params
	petID := c.Params("id")
	if petID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Pet ID required"})
	}

	// Ensure the pet belongs to the user
	var ownerEmail string
	err := config.DB.QueryRow(`
		SELECT u.email FROM user_profiles u
		JOIN pet_profiles p ON p.id = ANY(u.pets)
		WHERE p.id=$1 AND u.email=$2`, petID, userEmail).Scan(&ownerEmail)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: You don't own this pet"})
	} else if err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Delete pet from database
	_, err = config.DB.Exec("DELETE FROM pet_profiles WHERE id=$1", petID)
	if err != nil {
		log.Println("Database error deleting pet:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error deleting pet"})
	}

	return c.JSON(fiber.Map{"message": "Pet deleted successfully"})
}

// AssignCollarToPet assigns a collar to a pet such as unauthorized access, invalid input, or database failures.
func AssignCollarToPet(c *fiber.Ctx) error {
	petID := c.Params("id") // Get pet UUID from URL
	userEmail, ok := c.Locals("user").(string)

	// Ensure user is authenticated
	if !ok || userEmail == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req AssignCollarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate if the pet belongs to the user
	var ownerEmail string
	err := config.DB.QueryRow(`SELECT user_id FROM user_pets WHERE pet_id = $1`, petID).Scan(&ownerEmail)
	if errors.Is(err, sql.ErrNoRows) || ownerEmail != userEmail {
		log.Println("Check pet / user link failed", err)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized to assign collar to this pet"})
	}

	// Validate if the tag exists and is unassigned
	var isRegistered bool
	err = config.DB.QueryRow(`SELECT registered FROM reunite_collars WHERE tag_id = $1`, req.TagID).Scan(&isRegistered)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Collar not found error", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collar not found"})
	} else if isRegistered {
		log.Println("Collar is already assigned", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Collar is already assigned"})
	}

	// Set all other collars assigned to this pet as inactive
	_, err = config.DB.Exec(`UPDATE reunite_collars SET active = FALSE WHERE tag_id IN (SELECT tag_id FROM pet_collars WHERE pet_id = $1)`, petID)
	if err != nil {
		log.Println("Failed to deactivate collars:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to deactivate previous collars"})
	}

	// Assign the collar to the pet
	_, err = config.DB.Exec(`INSERT INTO pet_collars (pet_id, collar_id) VALUES ($1, $2) ON CONFLICT (pet_id, collar_id) DO NOTHING`, petID, req.TagID)
	if err != nil {
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while assigning collar"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Collar assigned and registered successfully"})
}

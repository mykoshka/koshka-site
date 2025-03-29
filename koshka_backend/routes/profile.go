package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/koshka_backend/config"
	"log"
)

// UserProfileResponse defines the response structure
type UserProfileResponse struct {
	Email                  string        `json:"email"`
	Name                   string        `json:"name"`
	Address                string        `json:"address"`
	MobileNumber           string        `json:"mobile_number"`
	JoinedOn               string        `json:"joined_on"`
	ProductPurchaseHistory []string      `json:"purchase_history"`
	ReuniteCollars         []string      `json:"reunite_collars"`
	Pets                   []PetResponse `json:"pets"`
}

// PetResponse defines pet details in user profile
type PetResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Picture     string   `json:"picture"`
	TagIDs      []string `json:"tag_ids"`
	CityLicence string   `json:"city_licence"`
	Neutered    bool     `json:"neutered"`
	Vaccinated  bool     `json:"vaccinated"`
	DOB         string   `json:"date_of_birth"`
}

// GetUserProfile fetches the user profile
func GetUserProfile(c *fiber.Ctx) error {
	userEmail := c.Locals("user").(string) // Extract from JWT

	var profile UserProfileResponse

	// Fetch user profile
	err := config.DB.QueryRow(`
        SELECT u.email, up.name, up.address, up.mobile_number, TO_CHAR(u.Created_at, 'YYYY-MM-DD')
        FROM users u 
        JOIN user_profiles up ON (u.email = up.email) WHERE u.email = $1`, userEmail).
		Scan(&profile.Email, &profile.Name, &profile.Address, &profile.MobileNumber, &profile.JoinedOn)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User profile not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while fetching user profile"})
	}

	// Fetch purchase history
	rows, err := config.DB.Query(`
        SELECT product_sku FROM purchase_history WHERE tag_id IN (
            SELECT tag_id FROM reunite_collars WHERE tag_id IN (
                SELECT unnest(reunite_collars) FROM user_profiles WHERE email = $1
            )
        )`, userEmail)

	if err == nil {
		for rows.Next() {
			var sku string
			if err := rows.Scan(&sku); err == nil {
				profile.ProductPurchaseHistory = append(profile.ProductPurchaseHistory, sku)
			}
		}
		rows.Close()
	}

	// Fetch reunite collars
	rows, err = config.DB.Query(`
        SELECT tag_id FROM reunite_collars AS rc JOIN public.pet_collars AS pc ON (rc.tag_id = pc.collar_id)
                      WHERE rc."active" = true AND pc.user_id = $1`, userEmail)

	if err == nil {
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err == nil {
				profile.ReuniteCollars = append(profile.ReuniteCollars, tag)
			}
		}
		rows.Close()
	}

	// Fetch pets
	rows, err = config.DB.Query(`
SELECT p.id, p.name, i.file, p.vaccinated, 
       p.neutered, TO_CHAR(p.date_of_birth, 'YYYY-MM-DD') AS date_of_birth, p.city_licence,
       array_agg(pc.collar_id) AS tag_ids
        FROM pet_profiles p
        LEFT JOIN user_pets up ON up.pet_id = p.id
        LEFT JOIN images i ON i.id = p.picture_id
        LEFT JOIN pet_collars pc ON pc.pet_id = p.id
        WHERE up.user_id = $1
        GROUP BY p.id, i.file`, userEmail)

	if err == nil {
		for rows.Next() {
			var pet PetResponse
			var tagIDs sql.NullString
			if err := rows.Scan(&pet.ID, &pet.Name, &pet.Picture, &pet.Vaccinated, &pet.Neutered, &pet.DOB, &pet.CityLicence, &tagIDs); err == nil {
				if tagIDs.Valid {
					pet.TagIDs = append(pet.TagIDs, tagIDs.String)
				}
				profile.Pets = append(profile.Pets, pet)
			}
		}
		rows.Close()
	} else {
		log.Println(err)
	}
	return c.Status(fiber.StatusOK).JSON(profile)
}

// UpdateUserProfile Updages the user profile
func UpdateUserProfile(c *fiber.Ctx) error {
	userEmail := c.Locals("user").(string) // Extract from JWT

	type ProfileUpdate struct {
		Name         string `json:"name"`
		Address      string `json:"address"`
		MobileNumber string `json:"mobile_number"`
	}

	var profile ProfileUpdate
	err := c.BodyParser(&profile)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update user profile
	_, err = config.DB.Exec(`UPDATE user_profiles SET name = $1, address = $2, mobile_number = $3 WHERE email = $4`,
		profile.Name, profile.Address, profile.MobileNumber, userEmail)

	if err != nil {
		log.Println("Error updating ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while updating user profile"})
	}

	return c.Status(fiber.StatusOK).JSON("Success")
}

// AddCollarToProfile Add Collar to profile
func AddCollarToProfile(c *fiber.Ctx) error {
	userEmail := c.Locals("user").(string) // Extract from JWT
	collarID := c.Params("tag")

	if collarID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Check if collar exists
	var count int
	err := config.DB.QueryRow(`SELECT COUNT(*) FROM reunite_collars WHERE tag_id = $1 AND 
                                           registered = false AND 
                                           active = true`, collarID).Scan(&count)
	if err == sql.ErrNoRows || count == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Request"})
	}

	// Add collar to profile
	_, err = config.DB.Exec(`INSERT INTO pet_collars (user_id, collar_id) VALUES ($1, $2)`, userEmail, collarID)
	if err != nil {
		log.Println("Error inserting ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while adding collar to profile"})
	}

	// Mark collar as assigned and registered
	_, err = config.DB.Exec(`UPDATE reunite_collars SET registered = TRUE WHERE tag_id = $1`, collarID)
	if err != nil {
		log.Println("Failed to register:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update collar status"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Collar added to profile"})
}

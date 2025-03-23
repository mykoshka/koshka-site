package routes

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/koshka_backend/config"
)

type AddApp struct {
	Name string `json:"name"`
}

// Struct to represent a single privilege
type APIPermission struct {
	RoutePattern string `json:"route_pattern"`
	Method       string `json:"method"`
	Allowed      bool   `json:"allowed"`
}

type Keys struct {
	ID        string `json:"id"`
	Owner     string `json:"created_by"`
	AppName   string `json:"app"`
	Key       string `json:"api_key"`
	CreatedOn string `json:"created_at"`
}

// BulkUploadTags - Admin-only route to bulk upload tags
func BulkUploadTags(c *fiber.Ctx) error {
	// ✅ Ensure only admins can access
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
	}

	// Parse uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File upload required"})
	}

	// Save file temporarily
	tempFilePath := fmt.Sprintf("uploads/%s", file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		log.Println("Error saving uploaded file:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Open CSV file
	csvFile, err := os.Open(tempFilePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer csvFile.Close()

	// Parse CSV
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Invalid CSV format"})
	}

	// Validate CSV headers
	if len(records) < 2 || records[0][0] != "TagID" || records[0][1] != "ProductSKU" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid CSV headers"})
	}

	// Process records
	insertCount := 0
	for _, row := range records[1:] {
		if len(row) < 4 {
			log.Println("Skipping incomplete row:", row)
			continue
		}

		tagID := row[0]
		productSKU := row[1]

		// Insert into database
		_, err := config.DB.Exec(`
			INSERT INTO reunite_collars (tag_id, product_sku)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (tag_id) DO UPDATE
			SET product_sku = EXCLUDED.product_sku`,
			tagID, productSKU)

		if err != nil {
			log.Println("Database error inserting tag:", err)
			continue
		}

		insertCount++
	}

	// Clean up temporary file
	os.Remove(tempFilePath)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Successfully uploaded %d tags", insertCount),
	})
}

// ListUsers - Admin-only route to fetch all users
func ListUsers(c *fiber.Ctx) error {
	// ✅ Ensure only admins can access
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
	}

	// Fetch all users from the database
	rows, err := config.DB.Query(`
		SELECT u.email, p.name, p.address, pr.name, u.locked, u.verified, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN user_profiles p ON u.email = p.email
		INNER JOIN personas pr ON u.persona_id = pr.id
		ORDER BY u.created_at DESC`)

	if err != nil {
		log.Println("Database error while fetching users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	defer rows.Close()

	// Store users in a slice
	users := []struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		Address   string `json:"address"`
		Persona   string `json:"persona"`
		Locked    string `json:"locked"`
		Verified  string `json:"verified"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{}

	for rows.Next() {
		var user struct {
			Email     string
			Name      sql.NullString
			Address   sql.NullString
			Persona   string
			Verified  bool
			Locked    bool
			CreatedAt string
			UpdatedAt string
		}

		if err := rows.Scan(&user.Email, &user.Name, &user.Address, &user.Persona, &user.Locked, &user.Verified, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println("Error scanning user:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}

		// Handle NULL values for name & address
		users = append(users, struct {
			Email     string `json:"email"`
			Name      string `json:"name"`
			Address   string `json:"address"`
			Persona   string `json:"persona"`
			Locked    string `json:"locked"`
			Verified  string `json:"verified"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}{
			Email:     user.Email,
			Name:      user.Name.String,
			Address:   user.Address.String,
			Persona:   user.Persona,
			Locked:    fmt.Sprintf("%t", user.Locked),
			Verified:  fmt.Sprintf("%t", user.Verified),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return c.JSON(users)
}

// ✅ List all registered routes (Admin only)
func ListAllRoutes(app *fiber.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ Ensure only admins can access
		isAdmin, ok := c.Locals("admin").(bool)
		if !ok || !isAdmin {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
		}

		// ✅ Allowed HTTP Methods
		var allowedMethods = map[string]bool{
			"GET":    true,
			"POST":   true,
			"PUT":    true,
			"DELETE": true,
		}

		// ✅ Map to group methods by path
		routeMap := make(map[string][]string)

		// ✅ Loop through all registered routes
		for _, route := range app.GetRoutes() {
			if allowedMethods[route.Method] { // ✅ Only store allowed methods
				routeMap[route.Path] = append(routeMap[route.Path], route.Method)
			}
		}

		// ✅ Convert map to sorted slice
		var routes []map[string]interface{}
		for path, methods := range routeMap {
			// ✅ Sort methods alphabetically
			sort.Strings(methods)

			routes = append(routes, map[string]interface{}{
				"path":    path,
				"methods": methods,
			})
		}

		// ✅ Sort routes by path
		sort.Slice(routes, func(i, j int) bool {
			return routes[i]["path"].(string) < routes[j]["path"].(string)
		})

		return c.JSON(fiber.Map{"routes": routes})
	}
}

// ✅ Create a New API Key
func CreateAPIKey(c *fiber.Ctx) error {
	userEmail := c.Locals("user").(string) // Extract from JWT
	// ✅ Ensure only admins can access
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
	}

	var req AddApp
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	var apiKey string
	err := config.DB.QueryRow("INSERT INTO api_keys (created_by, app) VALUES ($1, $2) RETURNING api_key", userEmail, req.Name).Scan(&apiKey)
	if err != nil {
		log.Println("error inserting in api_keys:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error with Database operation"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"api_key": apiKey})
}

// ✅ List All API Keys
func ListAPIKeys(c *fiber.Ctx) error {
	// ✅ Ensure only admins can access
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
	}

	var keys []Keys
	rows, err := config.DB.Query("SELECT id, created_by, app, api_key, created_at FROM api_keys")
	if err != nil {
		return err
	}
	defer rows.Close()

	//	var keys []map[string]interface{}
	for rows.Next() {
		var appKey Keys
		if err := rows.Scan(&appKey.ID, &appKey.Owner, &appKey.AppName, &appKey.Key, &appKey.CreatedOn); err != nil {
			log.Println("Error scanning keys:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		keys = append(keys, appKey)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"keys": keys})
}

// ✅ Grant or Update Route Permissions
func UpdateAPIPermission(c *fiber.Ctx) error {
	keyId := c.Params("key") // Get pet UUID from URL
	// ✅ Ensure only admins can access
	isAdmin, ok := c.Locals("admin").(bool)
	if !ok || !isAdmin {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: Admin access required"})
	}

	// Struct to hold multiple privileges
	type APIRequest struct {
		Privileges []APIPermission `json:"privileges"`
	}

	// Parse JSON request body
	var req APIRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		log.Println("❌ Error parsing JSON:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	// ✅ Validate input
	if len(req.Privileges) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No privileges provided"})
	}

	// ✅ Use a transaction for batch processing
	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("❌ Failed to start transaction:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	// ✅ Prepare SQL statement for inserting/updating permissions
	stmt, err := tx.Prepare(`
        INSERT INTO api_permissions (api_key_id, route_pattern, method, allowed)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (route_pattern, method)
        DO UPDATE SET allowed = EXCLUDED.allowed
    `)
	if err != nil {
		log.Println("❌ Error preparing SQL statement:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer stmt.Close()

	// ✅ Ensure method is always uppercase
	for _, p := range req.Privileges {
		methodUpper := strings.ToUpper(p.Method) // 🔹 Convert method to uppercase

		_, err := stmt.Exec(keyId, p.RoutePattern, methodUpper, p.Allowed)
		if err != nil {
			log.Println("❌ Error executing SQL:", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update permissions"})
		}
	}

	// ✅ Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("❌ Error committing transaction:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction failed"})
	}

	return c.JSON(fiber.Map{"message": "Permissions updated successfully"})
}

// ✅ Validate API Key & Permissions
func ValidateAPIKey(db *sql.DB, apiKey, route, method string) error {
	var apiKeyID int
	err := db.QueryRow("SELECT id FROM api_keys WHERE api_key = $1", apiKey).Scan(&apiKeyID)
	if err != nil {
		return errors.New("invalid API key")
	}

	var allowed bool
	err = db.QueryRow("SELECT allowed FROM api_permissions WHERE api_key_id = $1 AND route_pattern = $2 AND method = $3", apiKeyID, route, method).Scan(&allowed)
	if err != nil || !allowed {
		return errors.New("access denied")
	}

	return nil
}

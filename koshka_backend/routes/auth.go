package routes

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/koshka_backend/config"
	"github.com/koshka_backend/helpers"
	"github.com/koshka_backend/middleware"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"regexp"
	"strings"
)

// ChangePasswordRequest represents the request payload for changing the password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// Register represents the request payload for Registration
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Mobile   string `json:"mobile_number"`
}

// ✅ Generate a Secure Random Salt (16 bytes)
func GenerateSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Println("Error generating salt: ", err)
		panic("Failed to generate salt")
	}
	return base64.StdEncoding.EncodeToString(salt)
}

// ✅ Hash Password with Salt
func HashPassword(password, salt string) string {
	hashed := sha256.Sum256([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hashed[:])
}

// Login API
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	var storedPassword, salt string
	var locked bool

	// ✅ Retrieve Salt, Hashed & Locked Password from DB
	err := config.DB.QueryRow("SELECT password, salt, locked FROM users WHERE email = $1", input.Email).Scan(&storedPassword, &salt, &locked)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("Get pwd, salt, locked failed:", err)
		return c.Status(403).JSON(fiber.Map{"error": "Invalid credentials"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ Check if user is locked
	if locked {
		return c.Status(403).JSON(fiber.Map{"error": "Account locked. Contact support."})
	}

	// ✅ Hash Input Password with Stored Salt
	hashedInput := HashPassword(input.Password, salt)

	// ✅ Compare Hashed Passwords
	if hashedInput != storedPassword {
		log.Println("Authentication failed: ", hashedInput, storedPassword)
		return c.Status(403).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, _ := middleware.GenerateJWT(input.Email)

	// ✅ Update `updated_at` field after successful login
	_, err = config.DB.Exec("UPDATE users SET updated_at = NOW() WHERE email = $1", input.Email)
	if err != nil {
		log.Println("Database error updating updated_at:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error while updating login timestamp"})
	}

	// Store token in HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"message": "Login successful", "auth_token": token})
}

// Generate a random verification token
func generateVerificationToken() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Register API with Mandatory Field Validation and Verification URL
func Register(c *fiber.Ctx) error {

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check Mandatory Fields
	if input.Email == "" || input.Password == "" || input.Name == "" || input.Address == "" || input.Mobile == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Validate Email
	if !isValidEmail(input.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email"})
	}

	// Validate Password
	if !isValidPassword(input.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password does not meet minimum requirements"})
	}

	// ✅ Generate Unique Salt
	salt := GenerateSalt()

	// ✅ Hash the Password with Salt
	hashedPassword := HashPassword(input.Password, salt)

	// Generate Verification Token
	verificationToken := generateVerificationToken()
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", config.Config.BaseDomain, verificationToken)

	// ✅ Store User with Salt & Hashed Password
	_, err := config.DB.Exec("INSERT INTO users (email, password, salt, verification_token) VALUES ($1, $2, $3, $4)", input.Email, hashedPassword, salt, verificationToken)
	if err != nil {
		log.Println("Error creating user: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Insert Profile into DB
	_, err = config.DB.Exec("INSERT INTO user_profiles (email, name, address, mobile_number) VALUES ($1, $2, $3, $4)",
		input.Email, input.Name, input.Address, input.Mobile)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":         "Registration successful. Please verify your email.",
		"verificationURL": verificationURL,
	})
}

// Validate Email Format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Validate Password Complexity
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasNumber := strings.ContainsAny(password, "0123456789")
	return hasUpper && hasLower && hasNumber
}

// Verify Email API
func VerifyEmail(c *fiber.Ctx) error {
	verificationToken := c.Query("token")

	if verificationToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Verification token is required"})
	}

	// Check if token exists in the database
	var email string
	err := config.DB.QueryRow("SELECT email FROM users WHERE verification_token=$1 AND verified=false", verificationToken).Scan(&email)

	if err == sql.ErrNoRows {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired verification token"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Mark the user as verified and remove the token
	_, err = config.DB.Exec("UPDATE users SET verified=true, verification_token=NULL WHERE email=$1", email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Generate a JWT token for the verified user
	token, _ := middleware.GenerateJWT(email)

	// Store token in an HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Secure:   true,
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{"message": "Email verified successfully", "token": token})
}

// Request Password Reset API
func RequestPasswordReset(c *fiber.Ctx) error {
	type ResetRequest struct {
		Email string `json:"email"`
	}

	var input ResetRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if !isValidEmail(input.Email) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email"})
	}

	// Check if user exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", input.Email).Scan(&exists)
	if err != nil || !exists {
		return c.Status(400).JSON(fiber.Map{"error": "Email not found"})
	}

	// Generate a new verification token
	resetToken := generateVerificationToken()
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.Config.BaseDomain, resetToken)
	// Store the token in the database
	_, err = config.DB.Exec("UPDATE users SET verification_token=$1 WHERE email=$2", resetToken, input.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}
	templateID := os.Getenv("SENDGRID_TEMPLATE_PASSWORD_RESET")
	token := os.Getenv("SENDGRID_API_KEY")
	SendAuthEmail(input.Email, token, resetURL, templateID)
	return c.Status(200).JSON(fiber.Map{
		"message":  "Password reset link generated successfully.",
		"resetURL": resetURL,
	})
}

// Reset Password API
func ResetPassword(c *fiber.Ctx) error {
	type ResetInput struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var input ResetInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate token
	var email string
	err := config.DB.QueryRow("SELECT email FROM users WHERE verification_token=$1", input.Token).Scan(&email)

	if err == sql.ErrNoRows {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired token"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Validate Password
	if !isValidPassword(input.Password) {
		return c.Status(400).JSON(fiber.Map{"error": "Password does not meet minimum requirements"})
	}

	// Hash new password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	// ✅ Update password, unlock account, and clear reset token
	_, err = config.DB.Exec(`
		UPDATE users
		SET password = $1, locked = FALSE, verification_token = NULL, updated_at = NOW()
		WHERE email = $2`, hashedPassword, email)

	if err != nil {
		log.Println("Database error updating password:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Database error while updating password"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Password reset successful. You can now log in."})
}

// ChangePassword allows authenticated users to update their password
func ChangePassword(c *fiber.Ctx) error {
	userEmail := c.Locals("user").(string)

	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch current password hash from the database
	var storedHash string
	err := config.DB.QueryRow(`SELECT password FROM users WHERE email = $1`, userEmail).Scan(&storedHash)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while fetching user"})
	}

	// Validate old password
	if !helpers.CompareHash(req.OldPassword, storedHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Old password is incorrect"})
	}

	// Validate new password strength
	if !helpers.ValidatePassword(req.NewPassword) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password does not meet minimum requirements"})
	}

	// Hash the new password
	hashedPassword, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing new password"})
	}

	// Update password in database
	_, err = config.DB.Exec(`UPDATE users SET password = $1, updated_at = NOW() WHERE email = $2`, hashedPassword, userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while updating password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password changed successfully"})
}

// Send password reset email
func SendAuthEmail(email string, token string, resetURL string, templateID string) error {

	dynamicData := map[string]string{
		"subject": "Password Reset Request",
		"body":    fmt.Sprintf("Click <a href='%s'>here</a> to reset your password.", resetURL),
	}

	return helpers.SendEmailWithTemplate(email, "Password Reset Request", templateID, dynamicData)
}

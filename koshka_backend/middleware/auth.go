package middleware

import (
	"crypto/ecdsa"
	"errors"
	"github.com/koshka_backend/config"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtResponse struct {
	Token   string `json:"auth_token"`
	Refresh string `json:"refresh_token"`
}

var (
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

func addDays(days string, timeMesure string) int64 {
	var returnTime time.Duration
	i, err := strconv.ParseInt(days, 10, 64)
	if err != nil {
		return time.Time{}.Unix()
	}
	switch timeMesure {
	case "days":
		returnTime = time.Hour * (time.Duration(i * 24))
	case "hours":
		returnTime = time.Hour * (time.Duration(i))
	case "minutes":
		returnTime = time.Minute * (time.Duration(i))
	}

	return time.Now().Add(returnTime).Unix()
}

func createToken(sub string, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": sub,
		"exp": exp,
	})
	returnToken, err := token.SignedString(privateKey)
	return returnToken, err
}

func LoadJWTKeys() {
	privKeyPath := os.Getenv("JWT_PRIVATE_KEY")
	pubKeyPath := os.Getenv("JWT_PUBLIC_KEY")

	privKeyBytes, err := os.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}
	privKey, err := jwt.ParseECPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	privateKey = privKey

	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}
	pubKey, err := jwt.ParseECPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	publicKey = pubKey
}

func GenerateJWT(userID string) (JwtResponse, error) {
	var returnToken JwtResponse
	var err error
	if privateKey == nil {
		log.Println("Private key is not loaded")
		return returnToken, errors.New("private key not loaded")
	}

	returnToken.Token, err = createToken(userID, addDays(os.Getenv("JWT_TOKEN_DURATION"), "minutes"))
	returnToken.Refresh, err = createToken("1", addDays(os.Getenv("JWT_REFRESH_DURATION"), "days"))
	if err != nil {
		var emptyToken JwtResponse
		log.Printf("Error signing token: %v", err)
		return emptyToken, err
	}

	return returnToken, nil
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})
}

func revalidateJWT(tokenStr string, userid string) bool {
	return true
}

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		validationStr := c.Get("Refresh")
		if tokenStr == "" { // Get tokes from cookies instead of headers
			tokenStr = c.Cookies("auth_token")
			validationStr = c.Cookies("refresh_token")
			if tokenStr == "" {
				log.Printf("Missing token, error")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
			}
		}

		token, err := ValidateJWT(tokenStr)
		userid, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
		if err != nil || !token.Valid {
			_, err := ValidateJWT(validationStr)
			if err != nil {
				log.Printf("Token Validation error: %s", err)
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
			}

			newTokens, err := GenerateJWT(userid)
			// Store token in HTTP-only cookie
			c.Locals("new-auth", newTokens.Token)
			c.Locals("new-refresh", newTokens.Refresh)
		}

		if !ok {
			log.Printf("Claims error: %s", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		// Fetch user profile
		var isAdmin bool
		err = config.DB.QueryRow(`
        SELECT p.admin
        FROM users u JOIN personas p ON (u.persona_id = p.id) WHERE u.email = $1`, userid).
			Scan(&isAdmin)

		c.Locals("user", userid)
		c.Locals("admin", isAdmin)
		return c.Next()
	}
}

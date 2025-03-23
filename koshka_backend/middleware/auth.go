package middleware

import (
	"crypto/ecdsa"
	"errors"
	"github.com/koshka_backend/config"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

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

func GenerateJWT(userID string) (string, error) {
	if privateKey == nil {
		log.Println("Private key is not loaded")
		return "", errors.New("private key not loaded")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": userID,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
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

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		token, err := ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		// Fetch user profile
		var isAdmin bool
		err = config.DB.QueryRow(`
        SELECT p.admin
        FROM users u JOIN personas p ON (u.persona_id = p.id) WHERE u.email = $1`, claims["sub"]).
			Scan(&isAdmin)

		c.Locals("user", claims["sub"])
		c.Locals("admin", isAdmin)
		return c.Next()
	}
}

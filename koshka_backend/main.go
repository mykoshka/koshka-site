package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/koshka_backend/config"
	"github.com/koshka_backend/middleware"
	"github.com/koshka_backend/routes"
	"log"
	"time"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	app := fiber.New(fiber.Config{
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		AppName:           "Koshka Backend",
		ServerHeader:      "Koshka Backend Server",
		StrictRouting:     false,
	})

	// Set up caching middleware
	app.Use(cache.New(cache.Config{
		Expiration:   15 * time.Minute,
		CacheControl: true,
	}))

	// Set up CORS middleware
	app.Use(cors.New(cors.Config{
		AllowHeaders: "authorization, Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
	}))

	// Public Routes
	app.Post("/login", routes.Login)
	app.Post("/register", routes.Register)
	app.Post("/request-password-reset", routes.RequestPasswordReset)
	app.Post("/reset-password", routes.ResetPassword)
	app.Get("/verify-email", routes.VerifyEmail)
	app.Get("/collar/:tag_id", routes.FetchCollarInfo)
	app.Static("/images", "./uploads")

	// Protected routes (authenticated users)
	protected := app.Group("/api", middleware.JWTMiddleware())
	protected.Get("/profile", routes.GetUserProfile)
	protected.Post("/pet/add", routes.AddPet)
	protected.Put("/pet/:id", routes.UpdatePet)
	protected.Delete("/pet/:id", routes.DeletePet)
	protected.Post("/pet/:id/assign-collar", routes.AssignCollarToPet)
	protected.Post("/profile/auth/change-password", routes.ChangePassword)
	protected.Post("/profile/tag/:tag", routes.AddCollarToProfile)

	// Protected routes (Application APIs)
	protected.Get("/v1/reunite/ext/:code", routes.FetchPhoneNumberByCode) // server 2 server

	// Protected routes (Admin users)
	protected.Post("/admin/upload-tags", routes.BulkUploadTags)
	protected.Get("/admin/ListUsers", routes.ListUsers)
	protected.Post("/admin/api_keys", routes.CreateAPIKey)
	protected.Post("/admin/api_keys/permissions/:key", routes.UpdateAPIPermission)
	protected.Get("/admin/api_keys", routes.ListAPIKeys)
	protected.Get("/admin/routes", routes.ListAllRoutes(app))
	port := config.GetEnv("PORT", "3000")
	fmt.Println("Server running on port " + port)
	middleware.LoadJWTKeys()
	err := app.Listen("localhost:" + port)
	if err != nil {
		log.Println("Error starting:", err)
		return
	}
}

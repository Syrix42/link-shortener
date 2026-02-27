package main

import (
	"log"

	"github.com/Syrix42/link-shortener/internal/api"
	AuthHandler "github.com/Syrix42/link-shortener/internal/api/controllers/auth"
	"github.com/Syrix42/link-shortener/internal/config"
	"github.com/Syrix42/link-shortener/internal/infra/crypto"
	"github.com/Syrix42/link-shortener/internal/infra/database"
	"github.com/Syrix42/link-shortener/internal/infra/repositories"
	AuthService "github.com/Syrix42/link-shortener/internal/services/auth"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/Syrix42/link-shortener/swagger"
)

// @title Link Shortener API
// @version 1.0
// @description API for user authentication and link-shortening services.
// @BasePath /api/v1
func main() {

	// 1 Load configs from .env files from here
	dbconfig, err := config.LoadDBConfig("")
	if err != nil {
		log.Fatalf("failed to load db config: %v", err)
	}
	//2 DataBase connection is Established through the method bellow
	db, err := database.Connect(dbconfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	// 3) Build dependencies (repo -> service -> handler)

	hasher := crypto.NewBcryptHasher()
	//Repository of Aggregate User if you Intent to use it its already injected
	userRepo := repositories.NewUserRepository(db)

	authenticationService := AuthService.NewRegisterService(userRepo, hasher)
	authenticationHandler := AuthHandler.NewHandler(authenticationService)

	// 4) Create app + routes

	app := fiber.New()

	apiGroup := app.Group("/api")
	v1 := apiGroup.Group("/v1")

	api.AuthRoutes(v1, authenticationHandler)
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	appCfg := config.LoadAppConfig()
	log.Fatal(app.Listen(appCfg.ListenAddr()))
}

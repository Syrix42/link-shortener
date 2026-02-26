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
)

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
	Userrepo := repositories.NewUserRepository(db)

	AuthenticationService := AuthService.NewRegisterService(Userrepo, hasher)
	AuthenticationHandler := AuthHandler.NewHandler(AuthenticationService)

	// 4) Create app + routes

	app := fiber.New()

	apiGroup := app.Group("/api")
	v1 := apiGroup.Group("/v1")

	api.AuthRoutes(v1, AuthenticationHandler)

	appCfg := config.LoadAppConfig()
	log.Fatal(app.Listen(appCfg.ListenAddr()))
}

package main

import (
	"log"

	"github.com/Syrix42/link-shortener/internal/api"
	AuthHandler "github.com/Syrix42/link-shortener/internal/api/controllers/auth"
	"github.com/Syrix42/link-shortener/internal/config"
	crypto "github.com/Syrix42/link-shortener/internal/infra/crypto/hashing"
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
	AccessSecretConfig, err := config.LoadPrivateAccessJWTKey()
	if err != nil {
		log.Fatalf("Could not Load the JWT Access Secret :%v", err)
	}

	RefreshSecretConfig, err := config.LoadPrivateRefreshJWTKey()

	if err != nil {
		log.Fatalf("Could not Load the JWT Refresh Secret :%v", err)
	}

	//3 DataBase connection is Established through the method bellow
	db, err := database.Connect(dbconfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	// 4) Build dependencies (repo -> service -> handler)

	hasher := crypto.NewBcryptHasher()
	Comparer := crypto.NewBcryptComparer()
	//5 Repository of Aggregate User if you Intent to use it its already Injected
	Userrepo := repositories.NewUserRepository(db)
	//6 Repository of Aggregate Session if you Intent to use it its already Injected
	SessionRepo := repositories.NewSessionRepository(db)
	RegisterationService := AuthService.NewRegisterService(Userrepo, hasher)
	LoginService := AuthService.NewLoginService(Userrepo, Comparer, SessionRepo, SessionRepo, RefreshSecretConfig, AccessSecretConfig)
	AuthenticationHandler := AuthHandler.NewHandler(RegisterationService, LoginService)

	// 7) Create app + routes

	app := fiber.New()

	apiGroup := app.Group("/api")
	v1 := apiGroup.Group("/v1")

	api.AuthRoutes(v1, AuthenticationHandler)

	appCfg := config.LoadAppConfig()
	log.Fatal(app.Listen(appCfg.ListenAddr()))
}

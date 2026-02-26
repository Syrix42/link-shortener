package api

import (
	"github.com/Syrix42/link-shortener/internal/api/controllers/auth"
	"github.com/gofiber/fiber/v2"
)

// General Routher
func AuthRoutes(v1 fiber.Router, h auth.Handlers) {
	g := v1.Group("/auth")
	g.Post("/register", h.Register)
	//g.Post("/login", h.Login)
	//g.Post("/revoke", h.Revoke)
	//g.Post("/logout", h.Logout)

}

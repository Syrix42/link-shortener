package auth

// Login a new user
//@ Summary , Authenticates a New user
//@Description Creates a new session inside server
//@Tags auth
//@ Accept json
//@ Produces json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
import (
	"errors"

	"github.com/Syrix42/link-shortener/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

func (r *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	ctx := c.UserContext()
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_json",
		})

	}

	AccessToken, RefreshToken, err := r.LoginService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidEmailFormat):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid_email_format",
			})
		case errors.Is(err, auth.ErrUserNotFound):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user_not_found",
			})
		case errors.Is(err, auth.ErrInvalidPassword):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_password",
			})

		case errors.Is(err, auth.ErrTooManyActiveSessions):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "too_many_active_sessions",
			})
		default:
			return c.SendStatus(fiber.StatusInternalServerError)

		}
	}
	resp := LoginResponse{Status: "Authenticated", JWTRefreshToken: RefreshToken, JWTAccessToken: AccessToken}

	return c.Status(fiber.StatusOK).JSON(resp)

}

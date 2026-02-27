package auth

import (
	"errors"

	"github.com/Syrix42/link-shortener/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

// Register a new user.
// @Summary Register a new user
// @Description Creates a user account with email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (r *Handler) Register(c *fiber.Ctx) error {

	var request RegisterRequest
	ctx := c.UserContext()
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_json",
		})
	}

	err := r.RegisterService.Register(ctx, request.Email, request.Password)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidEmailFormat):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid_email_format",
			})
		case errors.Is(err, auth.ErrUserAlreadyExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user_already_exists",
			})
		default:
			return c.SendStatus(fiber.StatusInternalServerError)

		}
	}
	resp := RegisterResponse{Status: "Accepted", Message: "User Registered Successfully"}

	return c.Status(fiber.StatusOK).JSON(resp)
}

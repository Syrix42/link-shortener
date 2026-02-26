package auth

import (
	"errors"

	"github.com/Syrix42/link-shortener/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

func (r *Handler) Register(c *fiber.Ctx) error {

	var request RegisterationRequest
	ctx := c.UserContext()
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid_json",
		})
	}

	err := r.Regiserservice.Register(ctx, request.Email, request.Password)

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
	resp := RegisterationResponce{Status: "Accepted", Messsege: "User Registerd Successfuly"}
	return c.Status(fiber.StatusOK).JSON(resp)

}

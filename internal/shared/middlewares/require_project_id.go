package middlewares

import (
	"context"
	"platform/internal/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequireProjectID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headerValue := c.Get("X-Project-ID")
		if headerValue == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error_message": "Missing project identifier in request header",
			})
		}

		projectID, err := uuid.Parse(headerValue)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error_message": "Invalid project identifier in request header",
			})
		}

		ctx := context.WithValue(c.UserContext(), shared.ProjectIDContextKey, projectID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

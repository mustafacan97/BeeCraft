package middlewares

import (
	"context"
	"platform/internal/shared"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ProjectIDInjector() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headerValue := c.Get("X-Project-ID")
		if headerValue != "" {
			if projectID, err := uuid.Parse(headerValue); err == nil {
				ctx := context.WithValue(c.UserContext(), shared.ProjectIDContextKey, projectID)
				c.SetUserContext(ctx)
			}
		}
		return c.Next()
	}
}

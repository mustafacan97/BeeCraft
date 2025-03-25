package handlers

import (
	"context"
	"errors"
	"platform/internal/application/validators"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	validation = validator.New()
)

func init() {
	validation.RegisterValidation("password", validators.PasswordValidator)
}

type Handler[I Request, O any] interface {
	Handle(ctx context.Context, req *I) (*Response[O], error)
}

func Serve[I, O any](h Handler[I, O]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req I

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			// TODO
		}

		if err := c.ParamsParser(&req); err != nil {
			// TODO
		}

		if err := c.QueryParser(&req); err != nil {
			// TODO
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			// TODO
		}

		if err := validation.Struct(&req); err != nil {
			// Extract validation errors
			validationErrors := err.(validator.ValidationErrors)
			var errorMessages []string

			for _, fieldError := range validationErrors {
				// For each validation error, you can handle it here and send a custom error message
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, fieldError.Field()+" is required")
				case "min":
					errorMessages = append(errorMessages, fieldError.Field()+" must have at least "+fieldError.Param()+" characters")
				case "max":
					errorMessages = append(errorMessages, fieldError.Field()+" must have max "+fieldError.Param()+" characters")
				case "email":
					errorMessages = append(errorMessages, "Please enter a valid email address")
				case "password":
					errorMessages = append(errorMessages, "Password must have at least a uppercase-lowercase and a numeric characters")
				}
			}

			// If there are validation errors, return them as a list with 400 Bad Request
			if len(errorMessages) > 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
			}
		}

		resp, err := h.Handle(c.UserContext(), &req)
		if err != nil {
			zap.L().Error("An error occurred during request handling", zap.Error(err))
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "An unexpected error occurred. Please try again later."})
		}

		err = c.Next()
		if fiberErr, ok := err.(*fiber.Error); ok && fiberErr != nil && fiberErr.Code != fiber.StatusNotFound {
			zap.L().Error("An error occurred during request processing", zap.Int("status", fiberErr.Code), zap.Error(fiberErr))
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": "An unexpected error occurred. Please try again later."})
		}

		return c.Status(resp.ResponseStatus).JSON(resp)
	}
}

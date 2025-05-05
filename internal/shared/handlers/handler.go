package handlers

import (
	"context"
	"errors"
	"platform/internal/shared/validators"

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

		ctx := c.UserContext()

		if err := c.ParamsParser(&req); err != nil {
			// TODO: log here
			return fiber.NewError(fiber.StatusBadRequest, "Invalid URL parameters")
		}

		if err := c.QueryParser(&req); err != nil {
			// TODO: log here
			return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
		}

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			// TODO: log here
			return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
		}

		// For validation, I started with Fiber middleware.
		// If needed in the future, I can move it to Mediator pipeline.
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
				case "gt":
					errorMessages = append(errorMessages, fieldError.Field()+" must be greater than "+fieldError.Param())
				case "gte":
					errorMessages = append(errorMessages, fieldError.Field()+" must be greater than or equal to "+fieldError.Param())
				case "lt":
					errorMessages = append(errorMessages, fieldError.Field()+" must be less than "+fieldError.Param())
				case "lte":
					errorMessages = append(errorMessages, fieldError.Field()+" must be less than or equal to "+fieldError.Param())
				}
			}

			// If there are validation errors, return them as a list with 400 Bad Request
			if len(errorMessages) > 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errorMessages})
			}
		}

		resp, err := h.Handle(ctx, &req)
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

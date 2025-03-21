package handler

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	validation = validator.New()
)

func init() {

}

type Request any
type Response any

type Handler[I Request, O Response] interface {
	// context propagation
	Handle(context.Context, *I) (*O, error)
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
			// TODO
		}

		resp, err := h.Handle(c.UserContext(), &req)
		if err != nil {
			// TODO
		}

		c.JSON(resp)

		err = c.Next()
		if fiberErr, ok := err.(*fiber.Error); ok && fiberErr != nil && fiberErr.Code != fiber.StatusNotFound {
			// TODO
		}

		return nil
	}
}

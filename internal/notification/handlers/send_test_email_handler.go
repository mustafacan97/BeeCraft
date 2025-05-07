package handlers

import (
	"context"
	"errors"
	"platform/internal/notification/mediatr/commands"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
)

type SendTestEmailRequest struct {
	ID    uuid.UUID `params:"id" json:"-" validate:"required,uuid4"`
	Email string    `params:"email" json:"-" validate:"required,email"`
}

type SendTestEmailHandler struct{}

func (h *SendTestEmailHandler) Handle(ctx context.Context, req *SendTestEmailRequest) (*baseHandler.Response[shared.HALResource], error) {
	query := queries.GetEmailAccountByIDQuery{ID: req.ID}
	resp, err := mediator.Send[*queries.GetEmailAccountByIDQuery, *queries.GetEmailAccountByIDQueryResponse](ctx, &query)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](err), nil
	}
	if resp == nil {
		return baseHandler.FailedResponse[shared.HALResource](errors.New("email account not found")), nil
	}

	command := commands.SendTestEmailCommand{
		EmailAccountID: req.ID,
		To:             req.Email,
	}
	_, err = mediator.Send[*commands.SendTestEmailCommand, *commands.SendTestEmailCommandResponse](ctx, &command)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](err), nil
	}

	return baseHandler.SuccessResponse(&shared.HALResource{}), nil
}

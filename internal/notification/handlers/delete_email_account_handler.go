package handlers

import (
	"context"
	"platform/internal/notification/commands"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
)

type DeleteEmailAccountRequest struct {
	ProjectID uuid.UUID `reqHeader:"X-Project-ID" json:"-" validate:"required,uuid4"`
	ID        uuid.UUID `params:"id" json:"-" validate:"required,uuid4"`
}

type DeleteEmailAccountHandler struct{}

func (h *DeleteEmailAccountHandler) Handle(ctx context.Context, req *DeleteEmailAccountRequest) (*baseHandler.Response[shared.HALResource], error) {
	// STEP-1
	ctx = context.WithValue(ctx, shared.ProjectIDContextKey, req.ProjectID)

	// STEP-2
	command := commands.DeleteEmailAccountCommand{ID: req.ID}
	_, commandErr := mediator.Send[*commands.DeleteEmailAccountCommand, *commands.DeleteEmailAccountCommandResponse](ctx, &command)
	if commandErr != nil {
		return baseHandler.FailedResponse[shared.HALResource](commandErr), nil
	}

	// STEP-3
	response := &shared.HALResource{
		Links: shared.HALLinks{
			"self": {
				Href:   "",
				Method: "GET",
				Title:  "Get Email Account",
			},
		},
	}
	return baseHandler.SuccessResponse(response), nil
}

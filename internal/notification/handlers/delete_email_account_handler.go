package handlers

import (
	"context"
	"platform/internal/notification/mediatr/commands"
	event_notification "platform/internal/notification/mediatr/notifications"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
)

type DeleteEmailAccountRequest struct {
	ProjectID uuid.UUID `reqHeader:"X-Project-ID" params:"-" json:"-" validate:"required,uuid"`
	Email     string    `reqHeader:"-"  params:"email" json:"-" validate:"required,email"`
}

type DeleteEmailAccountResponse struct {
}

type DeleteEmailAccountHandler struct{}

func (h *DeleteEmailAccountHandler) Handle(ctx context.Context, req *DeleteEmailAccountRequest) (*baseHandler.Response[DeleteEmailAccountResponse], error) {
	// STEP-1: Delete the email account
	command := commands.DeleteEmailAccountCommand{Email: req.Email}
	_, err := mediator.Send[*commands.DeleteEmailAccountCommand, *commands.DeleteEmailAccountCommandResponse](ctx, &command)
	if err != nil {
		return baseHandler.FailedResponse[DeleteEmailAccountResponse](err), nil
	}

	// STEP-2: Publish email account deleted notification
	notification := event_notification.NewEmailAccountDeletedEvent(req.ProjectID, req.Email)
	mediator.Publish(ctx, &notification)

	// STEP-3: Return hateoas links to client
	data := DeleteEmailAccountResponse{}
	response := baseHandler.SuccessResponse(&data)
	response.Links = hateoasLinksForDelete()
	return response, nil
}

func hateoasLinksForDelete() shared.HALLinks {
	return shared.HALLinks{
		"list": {
			Href:   "/v1/notification/email-accounts?p=1&ps=10",
			Method: "GET",
			Title:  "List all emails on the first page",
		},
	}
}

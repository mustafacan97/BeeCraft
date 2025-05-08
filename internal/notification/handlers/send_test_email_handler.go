package handlers

import (
	"context"
	"fmt"
	"platform/internal/notification/mediatr/commands"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
)

type SendTestEmailRequest struct {
	ProjectID uuid.UUID `reqHeader:"X-Project-ID" params:"-" query:"-" json:"-" validate:"required,uuid"`
	From      string    `reqHeader:"-" params:"from" query:"-" json:"-" validate:"required,email"`
	To        string    `reqHeader:"-" params:"-" query:"-" json:"to" validate:"required,email"`
}

type SendTestEmailResponse struct {
}

type SendTestEmailHandler struct{}

func (h *SendTestEmailHandler) Handle(ctx context.Context, req *SendTestEmailRequest) (*baseHandler.Response[SendTestEmailResponse], error) {
	// STEP-1: Get from email
	query := queries.GetEmailAccountByEmailQuery{Email: req.From}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return baseHandler.NotFoundResponse[SendTestEmailResponse](), nil
	}

	// STEP-2: Send test email command
	command := commands.SendTestEmailCommand{
		From: req.From,
		To:   req.To,
	}
	_, err = mediator.Send[*commands.SendTestEmailCommand, *commands.SendTestEmailCommandResponse](ctx, &command)
	if err != nil {
		return nil, err
	}

	// STEP-10: Return hateoas links to user
	respData := SendTestEmailResponse{}
	response := baseHandler.SuccessResponse(&respData)
	response.Links = hateoasLinksForTestEmail(req.From)
	return response, nil
}

func hateoasLinksForTestEmail(email string) shared.HALLinks {
	return shared.HALLinks{
		"delete": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "DELETE",
			Title:  "Delete this email account",
		},
		"list": {
			Href:   "/v1/notification/email-accounts?p=1&ps=10",
			Method: "GET",
			Title:  "List all emails on the first page",
		},
		"self": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "GET",
			Title:  "View this email account",
		},
		"update": {
			Href:   fmt.Sprintf("/v1/notification/email-accounts/%s", email),
			Method: "PUT",
			Title:  "Update this email account",
		},
	}
}

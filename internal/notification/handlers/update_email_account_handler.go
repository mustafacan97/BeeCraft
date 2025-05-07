package handlers

import (
	"context"
	"fmt"
	"platform/internal/notification/mediatr/commands"
	event_notification "platform/internal/notification/mediatr/notifications"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"

	"github.com/google/uuid"
)

type UpdateEmailAccountRequest struct {
	ProjectID    uuid.UUID `reqHeader:"X-Project-ID" params:"-" query:"-" json:"-" validate:"required,uuid4"`
	Email        string    `reqHeader:"-" params:"email" query:"-" json:"-" validate:"required,email"`
	DisplayName  string    `reqHeader:"-" params:"-" query:"-" json:"display_name" validate:"required,max=255"`
	Host         string    `reqHeader:"-" params:"-" query:"-" json:"host" validate:"required,hostname|ip,max=255"`
	Port         int       `reqHeader:"-" params:"-" query:"-" json:"port" validate:"required,min=1,max=65535"`
	EnableSSL    bool      `reqHeader:"-" params:"-" query:"-" json:"enable_ssl"`
	TypeID       int       `reqHeader:"-" params:"-" query:"-" json:"type_id" validate:"required,oneof=1 2 3"`
	Username     string    `reqHeader:"-" params:"-" query:"-" json:"username"`
	Password     string    `reqHeader:"-" params:"-" query:"-" json:"password"`
	ClientID     string    `reqHeader:"-" params:"-" query:"-" json:"client_id"`
	TenantID     string    `reqHeader:"-" params:"-" query:"-" json:"tenant_id"`
	ClientSecret string    `reqHeader:"-" params:"-" query:"-" json:"client_secret"`
}

type UpdateEmailAccountResponse struct {
}

type UpdateEmailAccountHandler struct{}

func (h *UpdateEmailAccountHandler) Handle(ctx context.Context, req *UpdateEmailAccountRequest) (*baseHandler.Response[UpdateEmailAccountResponse], error) {
	// STEP-1: Check if email already exists
	query := queries.GetEmailAccountByEmailQuery{Email: req.Email}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return baseHandler.NotFoundResponse[UpdateEmailAccountResponse](), nil
	}

	// Step-2: Create email account command
	command := commands.UpdateEmailAccountCommand{
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		Host:         req.Host,
		Port:         req.Port,
		EnableSSL:    req.EnableSSL,
		TypeID:       req.TypeID,
		Username:     req.Username,
		Password:     req.Password,
		ClientID:     req.ClientID,
		TenantID:     req.TenantID,
		ClientSecret: req.ClientSecret,
	}
	_, err = mediator.Send[*commands.UpdateEmailAccountCommand, *commands.UpdateEmailAccountCommandResponse](ctx, &command)
	if err != nil {
		return nil, err
	}

	// STEP-3: Publish email account deleted notification
	notification := event_notification.NewEmailAccountUpdatedEvent(req.ProjectID, req.Email)
	mediator.Publish(ctx, &notification)

	// STEP-4: Return hateoas links to user
	respData := UpdateEmailAccountResponse{}
	response := baseHandler.SuccessResponse(&respData)
	response.Links = hateoasLinksForUpdate(req.Email)
	return response, nil
}

func hateoasLinksForUpdate(email string) shared.HALLinks {
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

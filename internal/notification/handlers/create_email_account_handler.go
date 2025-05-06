package handlers

import (
	"context"
	"errors"
	"fmt"
	"platform/internal/notification/commands"
	"platform/internal/notification/queries"
	"platform/internal/shared"
	"platform/internal/shared/handlers"
	"platform/pkg/services/mediator"

	"github.com/google/uuid"
)

type CreateEmailAccountRequest struct {
	Email        string `json:"email" validate:"required,email"`
	DisplayName  string `json:"display_name" validate:"required,max=255"`
	Host         string `json:"host" validate:"required,hostname|ip,max=255"`
	Port         int    `json:"port" validate:"required,min=1,max=65535"`
	EnableSSL    bool   `json:"enable_ssl"`
	TypeID       int    `json:"type_id" validate:"required,oneof=1 2 3"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	TenantID     string `json:"tenant_id"`
	ClientSecret string `json:"client_secret"`
}

type CreateEmailAccountHandler struct{}

func (h *CreateEmailAccountHandler) Handle(ctx context.Context, req *CreateEmailAccountRequest) (*handlers.Response[shared.HALResource], error) {
	// STEP-1: Check if email account already exists
	query := queries.GetEmailAccountByEmailQuery{Email: req.Email}
	ea, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}
	if ea != nil {
		return handlers.ConflictResponse[shared.HALResource](errors.New("email account already registered")), nil
	}

	// STEP-2: Create a new email account
	command := commands.CreateEmailAccountCommand{
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
	resp, err := mediator.Send[*commands.CreateEmailAccountCommand, *commands.CreateEmailAccountCommandResponse](ctx, &command)
	if err != nil {
		return nil, err
	}

	// STEP-3: Returns hateoas links to user
	response := &shared.HALResource{
		Links: createHateoasLinks(resp.ID),
	}
	return handlers.CreatedResponse(response), nil
}

func createHateoasLinks(emailAccountID uuid.UUID) shared.HALLinks {
	return shared.HALLinks{
		"self": {
			Href:   fmt.Sprintf("/email-accounts/%s", emailAccountID),
			Method: "GET",
			Title:  "View this email account",
		},
		"update": {
			Href:   fmt.Sprintf("/email-accounts/%s", emailAccountID),
			Method: "PUT",
			Title:  "Update this email account",
		},
		"delete": {
			Href:   fmt.Sprintf("/email-accounts/%s", emailAccountID),
			Method: "DELETE",
			Title:  "Delete this email account",
		},
		"list": {
			Href:   "/email-accounts",
			Method: "GET",
			Title:  "List all email accounts",
		},
	}
}

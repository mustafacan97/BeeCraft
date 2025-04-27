package handlers

import (
	"context"
	"errors"
	"platform/internal/notification/commands"
	"platform/internal/notification/queries"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"
)

var (
	errEmailAccountAlreadyRegistered = errors.New("email account already registered")
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

func (h *CreateEmailAccountHandler) Handle(ctx context.Context, req *CreateEmailAccountRequest) (*baseHandler.Response[shared.HALResource], error) {
	// STEP-1: Check if email already exists
	query := queries.GetEmailAccountByEmailQuery{Email: req.Email}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](err), nil
	}
	if resp != nil {
		return baseHandler.FailedResponse[shared.HALResource](errEmailAccountAlreadyRegistered), nil
	}

	// Step-2: Create email account command
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
	_, commandErr := mediator.Send[*commands.CreateEmailAccountCommand, *commands.CreateEmailAccountCommandResponse](ctx, &command)
	if commandErr != nil {
		return baseHandler.FailedResponse[shared.HALResource](commandErr), nil
	}

	// Step-3: Return HATEOAS information
	response := &shared.HALResource{
		Links: shared.HALLinks{
			"self": {
				Href:   "",
				Method: "GET",
				Title:  "Get Email Account",
			},
		},
	}
	return baseHandler.CreatedResponse(response), nil
}

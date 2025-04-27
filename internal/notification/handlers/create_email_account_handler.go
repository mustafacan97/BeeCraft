package handlers

import (
	"context"
	"errors"
	"platform/internal/notification/commands"
	internalDomain "platform/internal/notification/domain"
	"platform/internal/notification/queries"
	"platform/internal/shared"
	"platform/pkg/services/mediator"

	baseHandler "platform/internal/shared/handlers"
)

var (
	errEmailAccountAlreadyRegistered = errors.New("email account already registered")
	errUnsupportedSMTPType           = errors.New("smtp type is not supported")
)

type CreateEmailAccountRequest struct {
	Email        string `json:"email" validate:"required,email"`
	DisplayName  string `json:"display_name" validate:"required,max=255"`
	Host         string `json:"host" validate:"required,hostname|ip,max=255"`
	Port         int    `json:"port" validate:"required,min=1,max=65535"`
	EnableSSL    bool   `json:"enable_ssl" validate:"required"`
	TypeID       int    `json:"type_id" validate:"required,oneof=1 2 3"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	TenantID     string `json:"tenant_id"`
	ClientSecret string `json:"client_secret"`
}

type CreateEmailAccountResponse struct {
	Links shared.HALLinks `json:"_links"`
}

type CreateEmailAccountHandler struct{}

func (h *CreateEmailAccountHandler) Handle(ctx context.Context, req *CreateEmailAccountRequest) (*baseHandler.Response[CreateEmailAccountResponse], error) {
	query := queries.GetEmailAccountByEmailQuery{Email: req.Email}
	resp, err := mediator.Send[*queries.GetEmailAccountByEmailQuery, *queries.GetEmailAccountByEmailQueryResponse](ctx, &query)
	if err != nil {
		return baseHandler.FailedResponse[CreateEmailAccountResponse](err), nil
	}
	if resp != nil {
		return baseHandler.FailedResponse[CreateEmailAccountResponse](errEmailAccountAlreadyRegistered), nil
	}

	command := commands.CreateEmailAccountCommand{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Host:        req.Host,
		Port:        req.Port,
		EnableSSL:   req.EnableSSL,
		TypeID:      req.TypeID,
	}

	if req.TypeID == internalDomain.Login {
		command.Username = req.Username
		command.Password = req.Password
	} else if req.TypeID == internalDomain.GmailOAuth2 {
		command.ClientID = req.ClientID
		command.ClientSecret = req.ClientSecret
	} else if req.TypeID == internalDomain.MicrosoftOAuth2 {
		command.ClientID = req.ClientID
		command.TenantID = req.TenantID
		command.ClientSecret = req.ClientSecret
	} else {
		return nil, errUnsupportedSMTPType
	}

	_, commandErr := mediator.Send[*commands.CreateEmailAccountCommand, *commands.CreateEmailAccountCommandResponse](ctx, &command)
	if commandErr != nil {
		return baseHandler.FailedResponse[CreateEmailAccountResponse](commandErr), nil
	}

	response := &CreateEmailAccountResponse{
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

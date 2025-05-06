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

var (
	errEmailAccountNotFound = errors.New("email account not found")
)

type UpdateEmailAccountRequest struct {
	ProjectID    uuid.UUID `reqHeader:"X-Project-ID" json:"-" validate:"required,uuid4"`
	ID           uuid.UUID `params:"id" json:"-" validate:"required,uuid4"`
	Email        string    `json:"email" validate:"required,email"`
	DisplayName  string    `json:"display_name" validate:"required,max=255"`
	Host         string    `json:"host" validate:"required,hostname|ip,max=255"`
	Port         int       `json:"port" validate:"required,min=1,max=65535"`
	EnableSSL    bool      `json:"enable_ssl"`
	TypeID       int       `json:"type_id" validate:"required,oneof=1 2 3"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	ClientID     string    `json:"client_id"`
	TenantID     string    `json:"tenant_id"`
	ClientSecret string    `json:"client_secret"`
}

type UpdateEmailAccountHandler struct{}

func (h *UpdateEmailAccountHandler) Handle(ctx context.Context, req *UpdateEmailAccountRequest) (*baseHandler.Response[shared.HALResource], error) {
	// STEP-1: Set the ProjectID to the context so that repositories can use it
	ctx = context.WithValue(ctx, shared.ProjectIDContextKey, req.ProjectID)

	// STEP-2: Check if email already exists
	query := queries.GetEmailAccountByIDQuery{ID: req.ID}
	resp, err := mediator.Send[*queries.GetEmailAccountByIDQuery, *queries.GetEmailAccountByIDQueryResponse](ctx, &query)
	if err != nil {
		return baseHandler.FailedResponse[shared.HALResource](err), nil
	}
	if resp == nil {
		return baseHandler.FailedResponse[shared.HALResource](errEmailAccountNotFound), nil
	}

	// Step-3: Create email account command
	command := commands.UpdateEmailAccountCommand{
		ID:           req.ID,
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
	_, commandErr := mediator.Send[*commands.UpdateEmailAccountCommand, *commands.UpdateEmailAccountCommandResponse](ctx, &command)
	if commandErr != nil {
		return baseHandler.FailedResponse[shared.HALResource](commandErr), nil
	}

	// Step-4: Return HATEOAS information
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

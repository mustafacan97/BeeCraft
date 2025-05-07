package commands

import (
	"context"
	"platform/internal/notification/repositories"
	voExternal "platform/pkg/domain/value_object"
)

type DeleteEmailAccountCommand struct {
	Email string
}

type DeleteEmailAccountCommandResponse struct {
}

type DeleteEmailAccountCommandHandler struct {
	repository repositories.EmailAccountRepository
}

func NewDeleteEmailAccountCommandHandler(repository repositories.EmailAccountRepository) *DeleteEmailAccountCommandHandler {
	return &DeleteEmailAccountCommandHandler{repository: repository}
}

func (c *DeleteEmailAccountCommandHandler) Handle(ctx context.Context, command *DeleteEmailAccountCommand) (*DeleteEmailAccountCommandResponse, error) {
	email, err := voExternal.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}
	err = c.repository.Delete(ctx, email)
	if err != nil {
		return nil, err
	}
	return &DeleteEmailAccountCommandResponse{}, nil
}

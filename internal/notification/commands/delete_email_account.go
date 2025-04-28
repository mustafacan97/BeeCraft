package commands

import (
	"context"
	"platform/internal/notification/repositories"

	"github.com/google/uuid"
)

type DeleteEmailAccountCommand struct {
	ID uuid.UUID
}

type DeleteEmailAccountCommandResponse struct {
}

type DeleteEmailAccountCommandHandler struct {
	repository repositories.EmailAccountRepository
}

func NewDeleteEmailAccountCommandHandler(repository repositories.EmailAccountRepository) DeleteEmailAccountCommandHandler {
	return DeleteEmailAccountCommandHandler{repository: repository}
}

func (c *DeleteEmailAccountCommandHandler) Handle(ctx context.Context, command *DeleteEmailAccountCommand) (*DeleteEmailAccountCommandResponse, error) {
	err := c.repository.Delete(ctx, command.ID)
	if err != nil {
		return nil, err
	}
	return &DeleteEmailAccountCommandResponse{}, nil
}

package queries

import (
	"context"
	"platform/internal/notification/repositories"
	"time"

	"github.com/google/uuid"
)

type ListEmailAccountsQuery struct {
	Page     int
	PageSize int
}

type ListEmailAccountsQueryResponse struct {
	TotalCount int
	Data       []emailAccountDataForListing
}

type emailAccountDataForListing struct {
	ID          uuid.UUID
	Email       string
	DisplayName string
	TypeId      int
	CreatedAt   time.Time
}

type ListEmailAccountsQueryHandler struct {
	repository repositories.EmailAccountRepository
}

func NewListEmailAccountQueryHandler(repository repositories.EmailAccountRepository) *ListEmailAccountsQueryHandler {
	return &ListEmailAccountsQueryHandler{repository: repository}
}

func (c *ListEmailAccountsQueryHandler) Handle(ctx context.Context, query *ListEmailAccountsQuery) (*ListEmailAccountsQueryResponse, error) {
	emailAccounts, totalCount, err := c.repository.GetAll(ctx, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	response := ListEmailAccountsQueryResponse{
		TotalCount: totalCount,
		Data:       make([]emailAccountDataForListing, 0, len(emailAccounts)),
	}

	for _, acc := range emailAccounts {
		response.Data = append(response.Data, emailAccountDataForListing{
			ID:          acc.ID,
			Email:       acc.GetEmail().GetValue(),
			DisplayName: acc.GetDisplayName(),
			TypeId:      acc.GetSmtpType(),
			CreatedAt:   acc.GetCreatedDate(),
		})
	}

	return &response, nil
}

package queries

import (
	"context"
	"platform/internal/notification/repositories"
	"time"
)

type GetAllEmailAccountQuery struct {
	Page     int
	PageSize int
}

type GetAllEmailAccountQueryResponse struct {
	TotalCount int
	List       []data
}

type data struct {
	Email       string
	DisplayName string
	TypeId      int
	CreatedAt   time.Time
}

type GetAllEmailAccountQueryHandler struct {
	repository repositories.EmailAccountRepository
}

func NewGetAllEmailAccountQueryHandler(repository repositories.EmailAccountRepository) *GetAllEmailAccountQueryHandler {
	return &GetAllEmailAccountQueryHandler{repository: repository}
}

func (c *GetAllEmailAccountQueryHandler) Handle(ctx context.Context, query *GetAllEmailAccountQuery) (*GetAllEmailAccountQueryResponse, error) {
	emailAccounts, err := c.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	total := len(emailAccounts)
	start := (query.Page - 1) * query.PageSize
	if start > total {
		start = total
	}
	end := start + query.PageSize
	if end > total {
		end = total
	}
	pagedAccounts := emailAccounts[start:end]

	response := GetAllEmailAccountQueryResponse{
		TotalCount: total,
		List:       make([]data, 0, len(pagedAccounts)),
	}

	for _, acc := range emailAccounts {
		response.List = append(response.List, data{
			Email:       acc.GetEmail().Value(),
			DisplayName: acc.GetDisplayName(),
			TypeId:      acc.GetSmtpType(),
			CreatedAt:   acc.GetCreatedAt(),
		})
	}

	return &response, nil
}

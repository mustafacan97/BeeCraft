package handlers

import (
	"context"
	"platform/internal/notification/mediatr/queries"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"
	"time"

	"github.com/google/uuid"
)

type ListEmailAccountRequest struct {
	Page     int `query:"p" validate:"gt=0"`
	PageSize int `query:"ps" validate:"gt=0,lte=100"`
}

type ListEmailAccountResponse struct {
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

type ListEmailAccountHandler struct{}

func (h *ListEmailAccountHandler) Handle(ctx context.Context, req *ListEmailAccountRequest) (*baseHandler.Response[ListEmailAccountResponse], error) {
	query := queries.ListEmailAccountsQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	resp, err := mediator.Send[*queries.ListEmailAccountsQuery, *queries.ListEmailAccountsQueryResponse](ctx, &query)
	if err != nil {
		return nil, err
	}

	response := ListEmailAccountResponse{
		TotalCount: resp.TotalCount,
		Data:       make([]emailAccountDataForListing, 0, len(resp.Data)),
	}

	for _, data := range resp.Data {
		response.Data = append(response.Data, emailAccountDataForListing{
			ID:          data.ID,
			Email:       data.Email,
			DisplayName: data.DisplayName,
			TypeId:      data.TypeId,
			CreatedAt:   data.CreatedAt,
		})
	}

	return baseHandler.SuccessResponse(&response), nil
}

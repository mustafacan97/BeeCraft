package handlers

import (
	"context"
	"platform/internal/notification/mediatr/queries"
	"platform/internal/shared"
	baseHandler "platform/internal/shared/handlers"
	"platform/pkg/services/mediator"
	"time"

	"github.com/google/uuid"
)

type GetAllEmailAccountRequest struct {
	ProjectID uuid.UUID `reqHeader:"X-Project-ID" params:"-" query:"-" json:"-" validate:"required,uuid"`
	Page      int       `reqHeader:"-" params:"-" query:"p" json:"-" validate:"gt=0"`
	PageSize  int       `reqHeader:"-" params:"-" query:"ps" json:"-" validate:"gt=0,lte=100"`
}

type GetAllEmailAccountResponse struct {
	TotalCount int
	List       []data
}

type data struct {
	Email       string
	DisplayName string
	TypeId      int
	CreatedAt   time.Time
}

type GetAllEmailAccountHandler struct{}

func (h *GetAllEmailAccountHandler) Handle(ctx context.Context, req *GetAllEmailAccountRequest) (*baseHandler.Response[GetAllEmailAccountResponse], error) {
	// STEP-1: Get all email accounts
	query := &queries.GetAllEmailAccountQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	resp, err := mediator.Send[*queries.GetAllEmailAccountQuery, *queries.GetAllEmailAccountQueryResponse](ctx, query)
	if err != nil {
		return nil, err
	}

	// STEP-2: Create response struct
	respData := GetAllEmailAccountResponse{
		TotalCount: resp.TotalCount,
		List:       make([]data, 0, len(resp.List)),
	}

	// STEP-3: Fill the response data
	for _, li := range resp.List {
		respData.List = append(respData.List, data{
			Email:       li.Email,
			DisplayName: li.DisplayName,
			TypeId:      li.TypeId,
			CreatedAt:   li.CreatedAt,
		})
	}

	// STEP-4: Return data and hateoas links to user
	response := baseHandler.SuccessResponse(&respData)
	response.Links = hateoasLinksForAll()
	return response, nil
}

func hateoasLinksForAll() shared.HALLinks {
	return shared.HALLinks{
		"self": {
			Href:   "/v1/notification/email-accounts/:email",
			Method: "GET",
			Title:  "View this email account",
		},
		"delete": {
			Href:   "/v1/notification/email-accounts/:email",
			Method: "DELETE",
			Title:  "Delete this email account",
		},
	}
}

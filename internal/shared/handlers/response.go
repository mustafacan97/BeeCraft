package handlers

import "platform/internal/shared"

type Response[T any] struct {
	ResponseStatus int    `json:"-"`
	ErrorMessage   string `json:"error_message,omitempty"`
	Message        string `json:"message,omitempty"`
	Data           *T     `json:"data,omitempty"`
	shared.HALResource
}

func SuccessResponse[T any](data *T) *Response[T] {
	return &Response[T]{
		ResponseStatus: 200,
		Data:           data,
	}
}

func NotFoundResponse[T any]() *Response[T] {
	return &Response[T]{
		ResponseStatus: 404,
	}
}

func FailedResponse[T any](err error) *Response[T] {
	return &Response[T]{
		ResponseStatus: 200,
		ErrorMessage:   err.Error(),
	}
}

func CreatedResponse[T any](data *T) *Response[T] {
	return &Response[T]{
		ResponseStatus: 201,
		Data:           data,
	}
}

func CreatedResponseWithoutData[T any]() *Response[T] {
	return &Response[T]{
		ResponseStatus: 201,
	}
}

func ConflictResponse[T any](err error) *Response[T] {
	return &Response[T]{
		ResponseStatus: 409,
		ErrorMessage:   err.Error(),
	}
}

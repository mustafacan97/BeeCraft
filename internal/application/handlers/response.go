package handlers

type Response[T any] struct {
	ResponseStatus int    `json:"-"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	Message        string `json:"message,omitempty"`
	Data           *T     `json:"data,omitempty"`
}

func SuccessResponse[T any](data *T) *Response[T] {
	return &Response[T]{
		ResponseStatus: 200,
		Data:           data,
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

package mediator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ahmetb/go-linq/v3"
)

var requestHandlersRegistrations = map[reflect.Type]interface{}{}

// In the cases we don't need a response from our request handler, we can use `Unit` type, that actually is an empty struct.
type Unit struct{}

type RequestHandler[TRequest any, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

type RequestHandlerFunc func(ctx context.Context) (interface{}, error)

type RequestHandlerFactory[TRequest any, TResponse any] func() RequestHandler[TRequest, TResponse]

func RegisterRequestHandler[TRequest any, TResponse any](handler RequestHandler[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](handler)
}

func RegisterRequestHandlerFactory[TRequest any, TResponse any](factory RequestHandlerFactory[TRequest, TResponse]) error {
	return registerRequestHandler[TRequest, TResponse](factory)
}

func ClearRequestRegistrations() {
	requestHandlersRegistrations = map[reflect.Type]interface{}{}
}

func Send[TRequest any, TResponse any](ctx context.Context, request TRequest) (TResponse, error) {
	requestType := reflect.TypeOf(request)
	handler, ok := requestHandlersRegistrations[requestType]
	if !ok {
		return *new(TResponse), fmt.Errorf("no handler for request %T", request)
	}

	handlerValue, ok := buildRequestHandler[TRequest, TResponse](handler)
	if !ok {
		return *new(TResponse), fmt.Errorf("handler for request %T is not a Handler", request)
	}

	if len(pipelineBehaviors) == 0 {
		res, err := handlerValue.Handle(ctx, request)
		if err != nil {
			return *new(TResponse), fmt.Errorf("error handling request: %v", err)
		}
		return res, nil
	}

	var reversPipes = reversOrder(pipelineBehaviors)
	var lastHandler RequestHandlerFunc = func(ctx context.Context) (interface{}, error) {
		return handlerValue.Handle(ctx, request)
	}

	aggregateResult := linq.From(reversPipes).AggregateWithSeedT(lastHandler, func(next RequestHandlerFunc, pipe PipelineBehavior) RequestHandlerFunc {
		pipeValue := pipe
		nexValue := next

		var handlerFunc RequestHandlerFunc = func(ctx context.Context) (interface{}, error) {
			return pipeValue.Handle(ctx, request, nexValue)
		}

		return handlerFunc
	})

	v := aggregateResult.(RequestHandlerFunc)
	res, err := v(ctx)
	if err != nil {
		return *new(TResponse), fmt.Errorf("error handling request: %v", err)
	}

	response, ok := res.(TResponse)
	if !ok {
		return *new(TResponse), fmt.Errorf("handler returned unexpected type, expected %T", *new(TResponse))
	}

	return response, nil
}

func registerRequestHandler[TRequest any, TResponse any](handler any) error {
	var request TRequest
	requestType := reflect.TypeOf(request)

	_, exist := requestHandlersRegistrations[requestType]
	if exist {
		// each request in request/response strategy should have just one handler
		return fmt.Errorf("registered handler already exists in the registry for message %s", requestType.String())
	}

	requestHandlersRegistrations[requestType] = handler

	return nil
}

func buildRequestHandler[TRequest any, TResponse any](handler any) (RequestHandler[TRequest, TResponse], bool) {
	handlerValue, ok := handler.(RequestHandler[TRequest, TResponse])
	if !ok {
		factory, ok := handler.(RequestHandlerFactory[TRequest, TResponse])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

func reversOrder(values []interface{}) []interface{} {
	var reverseValues []interface{}

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}

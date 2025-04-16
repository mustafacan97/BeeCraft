package mediator

import (
	"context"
	"errors"
	"reflect"
)

var pipelineBehaviors []interface{} = []interface{}{}

type PipelineBehavior interface {
	Handle(ctx context.Context, request interface{}, next RequestHandlerFunc) (interface{}, error)
}

func RegisterRequestPipelineBehaviors(behaviors ...PipelineBehavior) error {
	for _, behavior := range behaviors {
		behaviorType := reflect.TypeOf(behavior)

		existsPipe := existsPipeType(behaviorType)
		if existsPipe {
			return errors.New("registered behavior already exists in the registry")
		}

		pipelineBehaviors = append(pipelineBehaviors, behavior)
	}

	return nil
}

func existsPipeType(p reflect.Type) bool {
	for _, pipe := range pipelineBehaviors {
		if reflect.TypeOf(pipe) == p {
			return true
		}
	}

	return false
}

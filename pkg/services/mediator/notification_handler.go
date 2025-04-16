package mediator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var notificationHandlersRegistrations = map[reflect.Type][]interface{}{}

type NotificationHandler[TNotification any] interface {
	Handle(ctx context.Context, notification TNotification) error
}

type NotificationHandlerFactory[TNotification any] func() NotificationHandler[TNotification]

func registerNotificationHandler[TEvent any](handler any) error {
	var event TEvent
	eventType := reflect.TypeOf(event)

	handlers, exist := notificationHandlersRegistrations[eventType]
	if !exist {
		notificationHandlersRegistrations[eventType] = []interface{}{handler}
		return nil
	}

	notificationHandlersRegistrations[eventType] = append(handlers, handler)

	return nil
}

func RegisterNotificationHandler[TEvent any](handler NotificationHandler[TEvent]) error {
	return registerNotificationHandler[TEvent](handler)
}

func RegisterNotificationHandlerFactory[TEvent any](factory NotificationHandlerFactory[TEvent]) error {
	return registerNotificationHandler[TEvent](factory)
}

func RegisterNotificationHandlers[TEvent any](handlers ...NotificationHandler[TEvent]) error {
	if len(handlers) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range handlers {
		err := RegisterNotificationHandler(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterNotificationHandlersFactories[TEvent any](factories ...NotificationHandlerFactory[TEvent]) error {
	if len(factories) == 0 {
		return errors.New("no handlers provided")
	}

	for _, handler := range factories {
		err := RegisterNotificationHandlerFactory(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearNotificationRegistrations() {
	notificationHandlersRegistrations = map[reflect.Type][]interface{}{}
}

func buildNotificationHandler[TNotification any](handler any) (NotificationHandler[TNotification], bool) {
	handlerValue, ok := handler.(NotificationHandler[TNotification])
	if !ok {
		factory, ok := handler.(NotificationHandlerFactory[TNotification])
		if !ok {
			return nil, false
		}

		return factory(), true
	}

	return handlerValue, true
}

func Publish[TNotification any](ctx context.Context, notification TNotification) error {
	eventType := reflect.TypeOf(notification)

	handlers, ok := notificationHandlersRegistrations[eventType]
	if !ok {
		return nil
	}

	for _, handler := range handlers {
		handlerValue, ok := buildNotificationHandler[TNotification](handler)

		if !ok {
			return fmt.Errorf("handler for notification %T is not a Handler", notification)
		}

		err := handlerValue.Handle(ctx, notification)
		if err != nil {
			return fmt.Errorf("error handling notification: %v", err)
		}
	}

	return nil
}

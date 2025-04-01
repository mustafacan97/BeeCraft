package domain

import "github.com/google/uuid"

type AggregateRoot interface {
	GetId() uuid.UUID
}

type BaseAggregateRoot struct {
	Id uuid.UUID
}

func (a *BaseAggregateRoot) GetId() uuid.UUID { return a.Id }

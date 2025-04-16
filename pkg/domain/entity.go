package domain

type Entity[ID comparable] interface {
	GetID() ID
	Equals(other BaseEntity[ID]) bool
}

// With the comparable constraint, we guarantee types with
// operations like ==, !=. This ensures comparability of IDs.
type BaseEntity[ID comparable] struct {
	ID ID
}

func NewBaseEntityWithID[ID comparable](id ID) BaseEntity[ID] {
	return BaseEntity[ID]{ID: id}
}

func (e BaseEntity[ID]) GetID() ID {
	return e.ID
}

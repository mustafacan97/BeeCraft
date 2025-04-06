package domain

import "reflect"

type ValueObject interface {
	Equals(other ValueObject) bool
	GetAtomicValues() []interface{}
}

type BaseValueObject struct{}

func (v BaseValueObject) Equals(other ValueObject) bool {
	if other == nil {
		return false
	}
	return reflect.DeepEqual(v.GetAtomicValues(), other.GetAtomicValues())
}

func (v *BaseValueObject) GetAtomicValues() []interface{} {
	// We will override it according to each concrete value object.
	return nil
}

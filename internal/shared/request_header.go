package shared

import "github.com/google/uuid"

const HEADER_ALIAS = "header"

type SharedRequestHeader struct {
	ProjectID uuid.UUID `header:"X-Project-ID" validate:"required,uuid4"`
}

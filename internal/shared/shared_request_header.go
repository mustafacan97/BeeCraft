package shared

type ctxKey string

const (
	ProjectIDHeader            = "X-Project-ID"
	ProjectIDContextKey ctxKey = ctxKey(ProjectIDHeader)
)

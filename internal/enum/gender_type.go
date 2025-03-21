package enum

// GenderType represents the gender of a user.
type GenderType int

const (
	NotSpecified GenderType = iota

	Male

	Female
)

// String method to return the string representation of the GenderType.
func (g GenderType) String() string {
	return [...]string{"not_specified", "male", "female"}[g]
}

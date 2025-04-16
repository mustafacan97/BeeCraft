package shared

// Language represents a language with its properties.
type Language struct {
	id      int
	culture string
	name    string
	rtl     bool
}

func (l Language) GetCulture() string         { return l.culture }
func (l Language) GetName() string            { return l.name }
func (l Language) IsRtl() bool                { return l.rtl }
func (l Language) Equals(other Language) bool { return l.id == other.id }

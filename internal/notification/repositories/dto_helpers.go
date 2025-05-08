package repositories

import "time"

// ptrToString creates a pointer to s if non-nil; otherwise returns nil.
func ptrToString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// ptrToStringValue returns a *string if the value is non-empty.
func ptrToStringValue(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ptrToTime returns time value or zero.
func ptrToTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

// ptrToTimeValue returns *time.Time if non-zero.
func ptrToTimeValue(t time.Time) *time.Time {
	if !t.IsZero() {
		return &t
	}
	return nil
}

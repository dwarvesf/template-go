package errorer

// Errorer is a interface for custom error type using in source
type Errorer interface {
	StatusCode() int
	Error() string
}

// Package usecase encapsulates use-case orchestration and business rules for the hotel
// domain, acting as the coordinator between the Web layer and the underlying SOAP
// services, responsible for input validation, request mapping, and error definitions.
package usecase

// ValidationError represents an input validation failure (an expected business error,
// mapped to HTTP 400). Unlike business errors returned by the upstream GDS, it occurs
// before the request enters domain logic.
type ValidationError struct {
	// Field is the name of the offending field (may be empty).
	Field string
	// Message is the human-readable message.
	Message string
}

// Error implements the error interface, returning a message that includes the field name
// (if non-empty) and the human-readable message.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// NewValidationError constructs an input validation error.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

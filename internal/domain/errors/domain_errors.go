package errors

// ValidationError represents a domain validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// BusinessRuleError represents a violation of business rules
type BusinessRuleError struct {
	Rule    string
	Message string
}

func (e *BusinessRuleError) Error() string {
	return e.Message
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return "resource not found: " + e.Resource + " with ID " + e.ID
}

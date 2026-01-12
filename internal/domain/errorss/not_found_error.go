package errorss

// ValidationError represents a domain validation error

// BusinessRuleError represents a violation of business rules

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return "resource not found: " + e.Resource + " with ID " + e.ID
}

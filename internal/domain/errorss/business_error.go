package errorss

type BusinessRuleError struct {
	Rule    string
	Message string
}

func (e *BusinessRuleError) Error() string {
	return e.Message
}

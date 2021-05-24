package errors

const ERROR_UNEXPECTED_ARGUMENT = 2

type unexpected struct {
	s string
}

// New returns an error that formats as the given text.
func NewUnexpected(text string) error {
	return &unexpected{text}
}

func (e *unexpected) Error() string {
	return e.s
}

func (e *unexpected) Type() int {
	return ERROR_UNEXPECTED_ARGUMENT
}

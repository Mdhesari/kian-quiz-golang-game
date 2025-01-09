package richerror

type Kind uint

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindUnAthorized
	KindNotFound
	KindUnexpected
	KindConflict
)

type RichError struct {
	kind      Kind
	operation string
	message   string
	err       error
	meta      map[string]interface{}
}

func New(op string, msg string) RichError {
	return RichError{
		operation: op,
		message:   msg,
	}
}

func (r RichError) WithMessage(msg string) RichError {
	r.message = msg

	return r
}

func (r RichError) WithOperation(op string) RichError {
	r.operation = op

	return r
}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind

	return r
}

func (r RichError) WithErr(err error) RichError {
	r.err = err

	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta

	return r
}

func (r RichError) Error() string {
	return r.operation + " : " + r.message
}

func (r RichError) Message() string {
	return r.message
}

func (r RichError) Kind() Kind {
	return r.kind
}

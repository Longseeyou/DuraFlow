package custom_error

import "errors"

var (
	ErrCreateInvalidRequest = errors.New("invalid create request")
	ErrUpdateInvalidRequest = errors.New("at least one field is required")
	ErrIdempotency          = errors.New("Idempotency")
)

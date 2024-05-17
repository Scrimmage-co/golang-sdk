package scrimmage

import "errors"

var (
	ErrInvalidURLProtocol error = errors.New("service url must start with protocol")
	ErrAccountIsNotLinked error = errors.New("selected account is not linked")
	ErrForbidden          error = errors.New("service token is invalid")
)

type BadRequestError struct {
	StatusCode int      `json:"statusCode"`
	Messages   []string `json:"message"`
	Err        string   `json:"error"`
}

func (e *BadRequestError) Error() string {
	return e.Err
}

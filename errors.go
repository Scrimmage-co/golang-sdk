package scrimmage

import "errors"

var (
	ErrInvalidURLProtocol error = errors.New("service url must start with protocol")
	ErrStatusCodeIsNotOK  error = errors.New("returned status code is not ok")
	ErrAccountIsNotLinked error = errors.New("selected account is not linked")
)

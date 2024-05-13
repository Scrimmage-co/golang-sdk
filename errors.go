package main

import "errors"

var (
	ErrInvalidURLProtocol error = errors.New("service url must start with protocol")
)

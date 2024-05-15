package scrimmage

import "strings"

func validateURLProtocol(url string, secure bool) bool {
	expectedProtocol := "http://"
	if secure {
		expectedProtocol = "https://"
	}

	return strings.HasPrefix(url, expectedProtocol)
}

func GetPtrOf[T any](input T) *T {
	return &input
}

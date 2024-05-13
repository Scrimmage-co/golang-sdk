package main

import "strings"

func validateURLProtocol(url string, secure bool) bool {
	expectedProtocol := "http://"
	if secure {
		expectedProtocol = "https://"
	}

	return strings.HasPrefix(url, expectedProtocol)
}

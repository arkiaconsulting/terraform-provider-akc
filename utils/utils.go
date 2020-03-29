package utils

import (
	"net/http"
)

// ResponseWasNotFound NotFound
func ResponseWasNotFound(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

// ResponseWasStatusCode With status code
func ResponseWasStatusCode(resp *http.Response, statusCode int) bool {
	if resp.StatusCode == statusCode {
		return true
	}

	return false
}

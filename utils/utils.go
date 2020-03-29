package utils

import (
	"net/http"
)

// ResponseWasNotFound NotFound
func ResponseWasNotFound(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

// ResponseWasThrottled TooManyRequests
func ResponseWasThrottled(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusTooManyRequests)
}

// ResponseWasForbidden Forbidden
func ResponseWasForbidden(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusForbidden)
}

// ResponseWasUnauthorized NotFound
func ResponseWasUnauthorized(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusUnauthorized)
}

// ResponseWasStatusCode With status code
func ResponseWasStatusCode(resp *http.Response, statusCode int) bool {
	if resp.StatusCode == statusCode {
		return true
	}

	return false
}

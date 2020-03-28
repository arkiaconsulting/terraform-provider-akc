package utils

import (
	"net/http"
)

func ResponseWasNotFound(resp *http.Response) bool {
	return ResponseWasStatusCode(resp, http.StatusNotFound)
}

func ResponseWasStatusCode(resp *http.Response, statusCode int) bool {
	if resp.StatusCode == statusCode {
		return true
	}

	return false
}

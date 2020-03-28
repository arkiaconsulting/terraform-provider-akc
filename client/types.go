package client

import "fmt"

// KeyValueResponse represents a Key Value response
type KeyValueResponse struct {
	Key          string
	Label        string
	ContentType  string `json:"content_type"`
	Value        string
	LastModified string `json:"last_modified"`
	Tags         map[string]string
}

type keyValueError struct {
	message string
	key     string
}

func kVError(key string, message string) error {
	return &keyValueError{key: key, message: message}
}

func (e *keyValueError) Error() string {
	return fmt.Sprintf("KV error <%s>: %s", e.key, e.message)
}

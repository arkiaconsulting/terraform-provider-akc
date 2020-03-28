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
	label   string
	key     string
}

func kVError(label string, key string, message string) error {
	return &keyValueError{label: label, key: key, message: message}
}

func (e *keyValueError) Error() string {
	return fmt.Sprintf("KV error %s/%s: %s", e.label, e.key, e.message)
}

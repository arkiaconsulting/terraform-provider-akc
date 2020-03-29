package client

// KeyValueResponse represents a Key Value response
type KeyValueResponse struct {
	Key          string
	Label        string
	ContentType  string `json:"content_type"`
	Value        string
	LastModified string `json:"last_modified"`
	Tags         map[string]string
}

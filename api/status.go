package api

import "encoding/json"

// A StatusMsg is a JSON response to an API request containing
// information about success, failures, etc
type StatusMsg struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewStatus returns a new instance of a StatusMsg
func NewStatus() *StatusMsg {
	return new(StatusMsg)
}

// NewStatusString creates a new StatusMsg with the provided status and message,
// and returns its contents as a JSON string
func NewStatusString(status, message string) string {
	s := NewStatus()
	s.Set(status, message)
	return s.String()
}

// Set sets the contents of a StatusMsg
func (s *StatusMsg) Set(status, message string) {
	s.Status = status
	s.Message = message
}

// String serializes the contents of a StatusMsg as a JSON string
func (s *StatusMsg) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes)
}

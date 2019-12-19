package api

import "encoding/json"

type StatusMsg struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewStatus(status, message string) *StatusMsg {
	s := new(StatusMsg)
	s.Set(status, message)
	return s
}

func NewStatusString(status, message string) string {
	s := NewStatus(status, message)
	return s.String()
}

func (s *StatusMsg) Set(status, message string) {
	s.Status = status
	s.Message = message
}

func (s *StatusMsg) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes)
}

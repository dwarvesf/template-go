package v1

type response struct {
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

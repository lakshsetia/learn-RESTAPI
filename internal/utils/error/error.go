package error

type ErrorResponse struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}
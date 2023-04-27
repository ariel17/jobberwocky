package http

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}
package models

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

package model

type ErrorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func NewErrorResponse(status int, msg string, data any) *ErrorResponse {
	return &ErrorResponse{
		Status: status,
		Msg:    msg,
	}
}

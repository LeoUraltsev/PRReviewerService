package helper

type ErrorResponse struct {
	Error Err `json:"error"`
}

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: Err{
			Code:    code,
			Message: message,
		},
	}
}

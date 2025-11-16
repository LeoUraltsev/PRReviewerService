package err

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

func NotFoundError() *ErrorResponse {
	return NewErrorResponse("NOT_FOUND", "resource not found")
}

func InternalServerError() *ErrorResponse {
	return NewErrorResponse("INTERNAL_SERVER_ERROR", "internal server error")
}

func TeamExistsError() *ErrorResponse {
	return NewErrorResponse("TEAM_EXISTS", "team_name already exists")
}

func PRExistsError() *ErrorResponse {
	return NewErrorResponse("PR_EXISTS", "PR id already exists")
}

func PRMergedError() *ErrorResponse {
	return NewErrorResponse("PR_MERGED", "cannot reassign on merged PR")
}

func IncorrectDataError() *ErrorResponse {
	return NewErrorResponse("INCORRECT_DATA", "incorrect data")
}

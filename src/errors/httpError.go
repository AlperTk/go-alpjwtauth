package errors

type HttpError struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NotFound() HttpError {
	return HttpError{
		404,
		"Not Found",
		"Requested resource not found",
	}
}
func UnauthorizedError() HttpError {
	return HttpError{
		401,
		"Unauthorized",
		"You are not authorized to access this resource",
	}
}
func ForbiddenRequestError() *HttpError {
	return &HttpError{
		403,
		"Forbidden",
		"Not authorized",
	}
}

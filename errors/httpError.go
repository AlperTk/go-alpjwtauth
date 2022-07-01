package errors

import "time"

var timeFormat = "2006-01-02T15:04:05.000:00"

type HttpError struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
}

func MethodNotAllowed() HttpError {
	return HttpError{
		time.Now().Format(timeFormat),
		405,
		"Method Not Allowed",
	}
}

func NotFound() HttpError {
	return HttpError{
		time.Now().Format(timeFormat),
		404,
		"Not Found",
	}
}

func ForbiddenRequestError() *HttpError {
	return &HttpError{
		time.Now().Format(timeFormat),
		403,
		"Forbidden",
	}
}

func UnauthorizedError() HttpError {
	return HttpError{
		time.Now().Format(timeFormat),
		401,
		"Unauthorized",
	}
}

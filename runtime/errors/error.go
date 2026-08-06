package errors

type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func BadRequest(msg string) error {

	return &HTTPError{
		Status:  400,
		Message: msg,
	}
}

func NotFound(msg string) error {

	return &HTTPError{
		Status:  404,
		Message: msg,
	}
}

func Unauthorized(msg string) error {

	return &HTTPError{
		Status:  401,
		Message: msg,
	}
}

func Internal(msg string) error {

	return &HTTPError{
		Status:  500,
		Message: msg,
	}
}
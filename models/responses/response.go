package responses

// Response's Message Constants
const (
	MessageSuccess string = "Success"
	MessageError   string = "Error"
)

// Standard Main Response Structure
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// NewResponse
// creates a new Response's instance with the given parameters
func NewResponse[T any](code int, message string, data T) Response[T] {
	if message == "" {
		return NewResponseWithoutMessage(code, data)
	}
	return Response[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewResponseWithoutMessage
// creates a new Response's instance with the given parameters but without message
func NewResponseWithoutMessage[T any](code int, data T) Response[T] {
	return Response[T]{
		Code: code,
		Data: data,
	}
}

// NewResponseWithoutData
// creates a new Response's instance with the given parameters but without data
func NewResponseWithoutData(code int, message string) Response[any] {
	return Response[any]{
		Code:    code,
		Message: message,
	}
}

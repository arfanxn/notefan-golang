package responses

// Response Message Constants
const (
	StatusSuccess string = "Success"
	StatusError   string = "Error"
)

type Response struct {
	body map[string]any
}

func NewResponse() Response {
	return Response{
		body: make(map[string]any),
	}
}

// Set the response code / status code
func (response Response) Code(statusCode int) Response {
	response.body["code"] = statusCode
	return response
}

// Get the response code / status code
func (response Response) GetCode() int {
	return response.body["code"].(int)
}

// Set the response message with status error
func (response Response) Error(msg string) Response {
	response.body["message"] = msg
	response.body["status"] = StatusError
	return response
}

// Set the response message with status success
func (response Response) Success(msg string) Response {
	response.body["message"] = msg
	response.body["status"] = StatusSuccess
	return response
}

// Set Response Body Specific Key
func (response Response) Body(key string, value any) Response {
	response.body[key] = value
	return response
}

// Get response body
func (response Response) GetBody() map[string]any {
	return response.body
}

package response

type HttpResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func NewHttpResponse(code int) *HttpResponse {
	return &HttpResponse{
		Code: code,
	}
}

func NewHttpResponseWithMessage(code int, message string) *HttpResponse {
	return &HttpResponse{
		Code:    code,
		Message: message,
	}
}

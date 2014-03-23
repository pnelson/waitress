package waitress

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Error interface {
	http.Handler
	SetHeaders(http.ResponseWriter)
}

type ErrorResponse struct {
	Code    int    `json:"-"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type MethodNotAllowedResponse struct {
	*ErrorResponse
	Allowed []string
}

func BadRequest() Error {
	return &ErrorResponse{
		Code: 400,
		Name: "Bad Request",
		Message: "The request could not be understood by the server " +
			"due to malformed syntax.",
	}
}

func Unauthorized() Error {
	return &ErrorResponse{
		Code: 401,
		Name: "Unauthorized",
		Message: "Authentication is required and has failed or " +
			"not yet been provided.",
	}
}

func Forbidden() Error {
	return &ErrorResponse{
		Code: 403,
		Name: "Forbidden",
		Message: "The request was valid, but the server is refusing " +
			"to respond to it.",
	}
}

func NotFound() Error {
	return &ErrorResponse{
		Code: 404,
		Name: "Not Found",
		Message: "The requested resource could not be found but " +
			"may be available again in the future.",
	}
}

func MethodNotAllowed(allowed []string) Error {
	return &MethodNotAllowedResponse{
		Allowed: allowed,
		ErrorResponse: &ErrorResponse{
			Code:    405,
			Name:    "Method Not Allowed",
			Message: "The method specified is not allowed for the resource.",
		},
	}
}

func NotAcceptable() Error {
	return &ErrorResponse{
		Code: 406,
		Name: "Not Acceptable",
		Message: "The requested resource is only capable of generating " +
			"content not acceptable according to the Accept headers sent " +
			"in the request.",
	}
}

func RequestTimeout() Error {
	return &ErrorResponse{
		Code: 408,
		Name: "Request Timeout",
		Message: "The client did not produce a request within the time " +
			"that the server was prepared to wait.",
	}
}

func Conflict() Error {
	return &ErrorResponse{
		Code: 409,
		Name: "Conflict",
		Message: "The request could not be completed due to a conflict " +
			"with the current state of the resource.",
	}
}

func Gone() Error {
	return &ErrorResponse{
		Code:    410,
		Name:    "Gone",
		Message: "The requested resource is no longer available.",
	}
}

func LengthRequired() Error {
	return &ErrorResponse{
		Code: 411,
		Name: "Length Required",
		Message: "The server refuses to accept the request without a " +
			"defined Content-Length.",
	}
}

func PreconditionFailed() Error {
	return &ErrorResponse{
		Code: 412,
		Name: "Precondition Failed",
		Message: "The server does not meet one or more of the " +
			"preconditions given in the request.",
	}
}

func RequestEntityTooLarge() Error {
	return &ErrorResponse{
		Code: 413,
		Name: "Request Entity Too Large",
		Message: "The request is larger than the server is willing or " +
			"able to process.",
	}
}

func RequestURITooLong() Error {
	return &ErrorResponse{
		Code:    414,
		Name:    "Request-URI Too Long",
		Message: "The provided URI was too long for the server to process.",
	}
}

func UnsupportedMediaType() Error {
	return &ErrorResponse{
		Code: 415,
		Name: "Unsupported Media Type",
		Message: "The request entity has a media type which the server " +
			"or resource does not support.",
	}
}

func RequestedRangeNotSatisfiable() Error {
	return &ErrorResponse{
		Code:    416,
		Name:    "Requested Range Not Satisfiable",
		Message: "The server cannot provide the requested range.",
	}
}

func ExpectationFailed() Error {
	return &ErrorResponse{
		Code: 417,
		Name: "Expectation Failed",
		Message: "The server cannot meet the requirements of the " +
			"Expect request-header field.",
	}
}

func UnprocessableEntity() Error {
	return &ErrorResponse{
		Code: 422,
		Name: "Unprocessable Entity",
		Message: "The request was well-formed but was unable to be " +
			"followed due to semantic errors.",
	}
}

func TooManyRequests() Error {
	return &ErrorResponse{
		Code: 429,
		Name: "Too Many Requests",
		Message: "The user has sent too many requests in a given " +
			"amount of time.",
	}
}

func InternalServerError() Error {
	return &ErrorResponse{
		Code: 500,
		Name: "Internal Server Error",
		Message: "The server encountered an unexpected condition which " +
			"prevented it from fulfilling the request.",
	}
}

func NotImplemented() Error {
	return &ErrorResponse{
		Code: 501,
		Name: "Not Implemented",
		Message: "The server does not support the functionality required " +
			"to fulfill the request.",
	}
}

func BadGateway() Error {
	return &ErrorResponse{
		Code: 502,
		Name: "Bad Gateway",
		Message: "The server, while acting as a gateway or proxy, received " +
			"an invalid response from the upstream server it accessed.",
	}
}

func ServiceUnavailable() Error {
	return &ErrorResponse{
		Code: 503,
		Name: "Service Unavailable",
		Message: "The server is currently unable to handle the request " +
			"due to a temporary overloading or maintenance of the server.",
	}
}

func GatewayTimeout() Error {
	return &ErrorResponse{
		Code: 504,
		Name: "Gateway Timeout",
		Message: "The server, while acting as a gateway or proxy, did " +
			"not receive a timely response from the upstream server.",
	}
}

func (e *ErrorResponse) SetHeaders(w http.ResponseWriter) {}
func (e *ErrorResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	e.SetHeaders(w)
	w.WriteHeader(e.Code)
	w.Write(b)
}

func (e *MethodNotAllowedResponse) SetHeaders(w http.ResponseWriter) {
	allowed := strings.Join(e.Allowed, ", ")
	w.Header().Set("Allow", allowed)
}

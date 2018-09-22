package pgway

import "net/http"

type ApiException struct {
	Message   string
	Error     error
	ErrorCode ErrorCode
}

type ErrorCode struct {
	Code       int
	Message    string
	HttpStatus int
}

var ApiNotFound = ErrorCode{Code: 0, Message: "api not found", HttpStatus: http.StatusNotFound}
var InvalidParameters = ErrorCode{Code: 1, Message: "invalid paramter", HttpStatus: http.StatusBadRequest}
var InternalServerError = ErrorCode{Code: 2, Message: "sorry", HttpStatus: http.StatusInternalServerError}

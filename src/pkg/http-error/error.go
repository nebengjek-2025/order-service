package error

import "net/http"

// CommonError struct
type CommonErrorData struct {
	Code         int         `json:"code"`
	ResponseCode int         `json:"responseCode,omitempty"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

// BadRequest struct
type BadRequestData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewBadRequest
func NewBadRequest() BadRequestData {
	errObj := BadRequestData{}
	errObj.Message = "Bad Request"
	errObj.Code = http.StatusBadRequest

	return errObj
}

// NotFound struct
type NotFoundData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewNotFound() NotFoundData {
	errObj := NotFoundData{}
	errObj.Message = "NotFound"
	errObj.Code = http.StatusNotFound

	return errObj
}

// Unauthorized struct
type UnauthorizedData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewUnauthorized() UnauthorizedData {
	errObj := UnauthorizedData{}
	errObj.Message = "Unauthorized"
	errObj.Code = http.StatusUnauthorized

	return errObj
}

// Conflict struct
type ConflictData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewConflict() ConflictData {
	errObj := ConflictData{}
	errObj.Message = "Conflict"
	errObj.Code = http.StatusConflict

	return errObj
}

// InternalServerError struct
type InternalServerErrorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewInternalServerError() InternalServerErrorData {
	errObj := InternalServerErrorData{}
	errObj.Message = "Internal Server Error"
	errObj.Code = http.StatusInternalServerError

	return errObj
}

type ErrorString struct {
	code    int
	message string
}

func (e ErrorString) Code() int {
	return e.code
}

func (e ErrorString) Error() string {
	return e.message
}

func (e ErrorString) Message() string {
	return e.message
}

// BadRequest will throw if the given request-body or params is not valid
func BadRequest(msg string) error {
	return &ErrorString{
		code:    http.StatusBadRequest,
		message: msg,
	}
}

// NotFound will throw if the requested item is not exists
func NotFound(msg string) error {
	return &ErrorString{
		code:    http.StatusNotFound,
		message: msg,
	}
}

// Conflict will throw if the current action already exists
func Conflict(msg string) error {
	return &ErrorString{
		code:    http.StatusConflict,
		message: msg,
	}
}

// InternalServerError will throw if any the Internal Server Error happen,
// Database, Third Party etc.
func InternalServerError(msg string) error {
	return &ErrorString{
		code:    http.StatusInternalServerError,
		message: msg,
	}
}

func UnauthorizedError(msg string) error {
	return &ErrorString{
		code:    http.StatusUnauthorized,
		message: msg,
	}
}

func ForbiddenError(msg string) error {
	return &ErrorString{
		code:    http.StatusForbidden,
		message: msg,
	}
}

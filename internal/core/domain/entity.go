package domain

import "net/http"

type ResponseT[T any] struct {
	Code    string            `json:"code"`
	Message interface{}       `json:"message"`
	Data    T                 `json:"data,omitempty"`
	Header  map[string]string `json:"-"`
}

type AppError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Code           string `json:"code"`
	Description    string `json:"description"`
}

func NewAppBadRequestError(respCode string, desc string) AppError {
	return AppError{
		HTTPStatusCode: http.StatusBadRequest,
		Code:           respCode,
		Description:    desc,
	}
}

func NewAppInternalServerError(respCode string, desc string) AppError {
	return AppError{
		HTTPStatusCode: http.StatusInternalServerError,
		Code:           respCode,
		Description:    desc,
	}
}

func NewAppUnauthorizedError(respCode string, desc string) AppError {
	return AppError{
		HTTPStatusCode: http.StatusUnauthorized,
		Code:           respCode,
		Description:    desc,
	}
}

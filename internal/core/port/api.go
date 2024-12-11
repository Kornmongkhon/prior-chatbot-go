package port

import (
	"context"
	"net/http"
)

type MyServer interface {
	GetRepo() MyRepo
	GetContext() context.Context
	GetRequest() *http.Request
	BindRequest(interface{}) error
	ToResponse(code int, statusCode string, msg interface{}, data interface{}) error
}

type Response struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

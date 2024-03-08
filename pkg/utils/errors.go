package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func NewNotFoundError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func GinContextToGenericContext(ctx *gin.Context) context.Context {
	nCtx := context.Background()
	nCtx = context.WithValue(nCtx, "request", ctx.Request)
	return nCtx
}

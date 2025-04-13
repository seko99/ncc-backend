package domain

import (
	"errors"
	"fmt"
)

var (
	ErrExecute           = errors.New("exec error")
	ErrParameters        = errors.New("parameters error")
	ErrPaymentNotAllowed = errors.New("payment not allowed")
)

type ErrResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func HttpErrorResponse(err error, format string, a ...any) ErrResponse {
	return ErrResponse{
		Error:   err.Error(),
		Message: fmt.Sprintf(format, a...),
	}
}

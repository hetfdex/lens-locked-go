package model

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	StatusCode int
	Message    string
}

func MustNotBeEmptyErrorMessage(value string) string {
	message := fmt.Sprintf("%s must not be empty", value)

	return message
}

func NewInternalServerApiError(message string) *ApiError {
	return newApiError(http.StatusInternalServerError, message)
}

func NewNotFoundApiError(message string) *ApiError {
	return newApiError(http.StatusNotFound, message)
}

func NewForbiddenApiError(message string) *ApiError {
	return newApiError(http.StatusForbidden, message)
}

func NewBadRequestApiError(message string) *ApiError {
	return newApiError(http.StatusBadRequest, message)
}

func NewConflictApiError(message string) *ApiError {
	return newApiError(http.StatusConflict, message)
}

func newApiError(statusCode int, message string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
	}
}

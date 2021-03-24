package model

import "net/http"

type ApiError struct {
	StatusCode int
	Message    string
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

func newApiError(statusCode int, message string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
	}
}

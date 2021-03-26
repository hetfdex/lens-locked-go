package model

import (
	"fmt"
	"net/http"
)

const internalServerErrorMessage = "Something went wrong"
const notFoundErrorMessage = "Resource not found"
const forbiddenErrorMessage = "Access forbidden"
const badRequestErrorMessage = "Invalid data provided"
const conflictErrorMessage = "Resource already exists"

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Alert() *Alert {
	var alertLevel string
	var message string

	switch e.StatusCode {
	case http.StatusInternalServerError:
		alertLevel = alertLevelError
		message = internalServerErrorMessage
	case http.StatusNotFound:
		alertLevel = alertLevelWarning
		message = notFoundErrorMessage
	case http.StatusForbidden:
		alertLevel = alertLevelError
		message = forbiddenErrorMessage
	case http.StatusBadRequest:
		alertLevel = alertLevelWarning
		message = badRequestErrorMessage
	case http.StatusConflict:
		alertLevel = alertLevelWarning
		message = conflictErrorMessage
	}
	return newAlert(alertLevel, message)
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

func MustNotBeEmptyErrorMessage(value string) string {
	message := fmt.Sprintf("%s must not be empty", value)

	return message
}

func newApiError(statusCode int, message string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Message:    message,
	}
}

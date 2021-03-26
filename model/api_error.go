package model

import (
	"fmt"
	"net/http"
	"strings"
)

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Alert() *Alert {
	var alertLevel string

	switch e.StatusCode {
	case http.StatusInternalServerError:
		alertLevel = alertLevelError
	case http.StatusNotFound:
		alertLevel = alertLevelWarning
	case http.StatusForbidden:
		alertLevel = alertLevelError
	case http.StatusBadRequest:
		alertLevel = alertLevelWarning
	case http.StatusConflict:
		alertLevel = alertLevelWarning
	}
	return newAlert(alertLevel, strings.Title(e.Message))
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

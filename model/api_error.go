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

	if e.StatusCode == http.StatusInternalServerError {
		alertLevel = AlertLevelError
		message = internalServerErrorMessage
	} else if e.StatusCode == http.StatusNotFound {
		alertLevel = AlertLevelWarning
		message = notFoundErrorMessage
	} else if e.StatusCode == http.StatusForbidden {
		alertLevel = AlertLevelError
		message = forbiddenErrorMessage
	} else if e.StatusCode == http.StatusBadRequest {
		alertLevel = AlertLevelInfo
		message = badRequestErrorMessage
	} else if e.StatusCode == http.StatusConflict {
		alertLevel = AlertLevelWarning
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

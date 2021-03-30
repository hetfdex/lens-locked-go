package model

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const genericErrorMessage = "something went wrong"

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Alert() *Alert {
	var alertLevel string

	message := e.Message

	log.Println(e.StatusCode, e.Message)

	switch e.StatusCode {
	case http.StatusInternalServerError:
		alertLevel = alertLevelError
		message = genericErrorMessage
	case http.StatusNotFound:
		alertLevel = alertLevelWarning
	case http.StatusForbidden:
		alertLevel = alertLevelError
	case http.StatusBadRequest:
		alertLevel = alertLevelWarning
	case http.StatusConflict:
		alertLevel = alertLevelWarning
	}
	return &Alert{
		Level:   alertLevel,
		Message: strings.Title(message),
	}
}

func NewInternalServerApiError(message string) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
	}
}

func NewNotFoundApiError(message string) *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func NewForbiddenApiError(message string) *Error {
	return &Error{
		StatusCode: http.StatusForbidden,
		Message:    message,
	}
}

func NewBadRequestApiError(message string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func NewConflictApiError(message string) *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

func MustNotBeEmptyErrorMessage(value string) string {
	message := fmt.Sprintf("%s must not be empty", value)

	return message
}

package model

const AlertLevelError = "danger"
const AlertLevelWarning = "warning"
const AlertLevelInfo = "info"
const AlertLevelSuccess = "success"

type Alert struct {
	Level   string
	Message string
}

func NewErrorAlert(message string) *Alert {
	return newAlert(AlertLevelError, message)
}

func NewWarningAlert(message string) *Alert {
	return newAlert(AlertLevelWarning, message)
}

func NewInfoAlert(message string) *Alert {
	return newAlert(AlertLevelInfo, message)
}

func NewSuccessAlert(message string) *Alert {
	return newAlert(AlertLevelSuccess, message)
}

func newAlert(level string, message string) *Alert {
	return &Alert{
		Level:   level,
		Message: message,
	}
}

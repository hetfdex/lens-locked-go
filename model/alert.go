package model

const alertLevelError = "danger"
const alertLevelWarning = "warning"
const alertLevelInfo = "info"
const alertLevelSuccess = "success"

type Alert struct {
	Level   string
	Message string
}

func NewSuccessAlert(message string) *Alert {
	return newAlert(alertLevelSuccess, message)
}

func newAlert(level string, message string) *Alert {
	return &Alert{
		Level:   level,
		Message: message,
	}
}

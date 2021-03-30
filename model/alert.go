package model

const alertLevelError = "danger"
const alertLevelWarning = "warning"
const alertLevelSuccess = "success"

type Alert struct {
	Level   string
	Message string
}

func NewSuccessAlert(message string) *Alert {
	return &Alert{
		Level:   alertLevelSuccess,
		Message: message,
	}
}

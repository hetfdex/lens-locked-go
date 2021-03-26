package model

const alertLevelError = "danger"
const alertLevelWarning = "warning"

type Alert struct {
	Level   string
	Message string
}

func newAlert(level string, message string) *Alert {
	return &Alert{
		Level:   level,
		Message: message,
	}
}

package model

const AlertLevelError = "danger"
const AlertLevelWarning = "warning"
const AlertLevelInfo = "info"
const AlertLevelSuccess = "success"

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

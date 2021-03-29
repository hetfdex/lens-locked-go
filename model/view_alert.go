package model

const alertLevelError = "danger"
const alertLevelWarning = "warning"
const alertLevelSuccess = "success"

type AlertView struct {
	Level   string
	Message string
}

func NewSuccessAlert(message string) *AlertView {
	return &AlertView{
		Level:   alertLevelSuccess,
		Message: message,
	}
}

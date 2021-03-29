package model

const alertLevelError = "danger"
const alertLevelWarning = "warning"

type AlertView struct {
	Level   string
	Message string
}

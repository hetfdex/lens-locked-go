package model

type DataView struct {
	Alert *AlertView
	User  *User
	Data  interface{}
}

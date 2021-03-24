package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
)

const pepper = "6Sk65RHhGW7S4qnVPV7m"

func generateFromPassword(user *model.User) error {
	pw := []byte(user.Password + pepper)

	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = ""
	user.PasswordHash = string(hash)

	return nil
}

func compareHashAndPassword(user *model.User, password string) error {
	hash := []byte(user.PasswordHash)
	pw := []byte(password + pepper)

	return bcrypt.CompareHashAndPassword(hash, pw)
}

package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
)

const pepper = "6Sk65RHhGW7S4qnVPV7m"

func generateFromPassword(user *model.User) *model.ApiError {
	pw := []byte(user.Password + pepper)

	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	user.Password = ""
	user.PasswordHash = string(hash)

	return nil
}

func compareHashAndPassword(user *model.User, password string) *model.ApiError {
	hash := []byte(user.PasswordHash)
	pw := []byte(password + pepper)

	err := bcrypt.CompareHashAndPassword(hash, pw)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(err.Error())
		}
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

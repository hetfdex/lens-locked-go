package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
)

func generateFromPassword(user *model.User) *model.ApiError {
	pw := []byte(user.Password + config.Pepper)

	h, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	user.Password = ""
	user.PasswordHash = string(h)

	return nil
}

func compareHashAndPassword(user *model.User, password string) *model.ApiError {
	h := []byte(user.PasswordHash)
	pw := []byte(password + config.Pepper)

	err := bcrypt.CompareHashAndPassword(h, pw)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(err.Error())
		}
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func generateRememberHash(hasher *hash.Hasher, user *model.User) *model.ApiError {
	h, err := hasher.GenerateHash(user.Remember)

	if err != nil {
		return err
	}
	user.Remember = ""
	user.RememberHash = h

	return nil
}

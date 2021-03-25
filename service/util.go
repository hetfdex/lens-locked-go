package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/config"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"strings"
)

func generateFromPassword(password string) (string, *model.ApiError) {
	if password == "" {
		return "", model.NewInternalServerApiError("password must not be empty")
	}
	pw := []byte(password + config.Pepper)

	pwHash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(pwHash), nil
}

func compareHashAndPassword(passwordHash string, password string) *model.ApiError {

	if passwordHash == "" {
		return model.NewInternalServerApiError("passwordHash must not be empty")
	}

	if password == "" {
		return model.NewInternalServerApiError("password must not be empty")
	}
	pwHash := []byte(passwordHash)
	pw := []byte(password + config.Pepper)

	err := bcrypt.CompareHashAndPassword(pwHash, pw)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError("invalid password")
		}
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func generateToken() (string, *model.ApiError) {
	token, err := rand.GenerateTokenString()

	if err != nil {
		return "", err
	}
	return token, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(email)
}

package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/config"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"lens-locked-go/validator"
)

func generateFromPassword(password string) (string, *model.ApiError) {
	apiErr := validator.StringNotEmpty("password", password)

	if apiErr != nil {
		return "", apiErr
	}
	pw := []byte(password + config.Pepper)

	pwHash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(pwHash), nil
}

func compareHashAndPassword(passwordHash string, password string) *model.ApiError {
	apiErr := validator.StringNotEmpty("passwordHash", passwordHash)

	if apiErr != nil {
		return apiErr
	}
	apiErr = validator.StringNotEmpty("password", password)

	if apiErr != nil {
		return apiErr
	}
	pwHash := []byte(passwordHash)
	pw := []byte(password + config.Pepper)

	err := bcrypt.CompareHashAndPassword(pwHash, pw)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(err.Error())
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

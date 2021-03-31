package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/config"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"regexp"
	"strings"
)

const invalidEmailErrorMessage = "email address must have a valid format"
const invalidPasswordErrorMessage = "invalid password"
const invalidPasswordLengthErrorMessage = "password must be at least 8 characters"

var emailRegex = regexp.MustCompile(`^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$`)

func lower(s string) string {
	return strings.ToLower(s)
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func generateHashFromPassword(password string) (string, *model.Error) {
	if password == "" {
		return "", model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("password"))
	}

	hs, err := bcrypt.GenerateFromPassword([]byte(password+config.Pepper), bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(hs), nil
}

func compareHashAndPassword(hash string, password string) *model.Error {

	if hash == "" {
		return model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("hash"))
	}

	if password == "" {
		return model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("password"))
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+config.Pepper))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(invalidPasswordErrorMessage)
		}
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func generateToken() (string, *model.Error) {
	token, err := rand.GenerateTokenString()

	if err != nil {
		return "", err
	}
	return token, nil
}

func validEmail(email string) *model.Error {
	if len(email) < 3 && len(email) > 254 {
		return model.NewBadRequestApiError(invalidEmailErrorMessage)
	}

	if !emailRegex.MatchString(email) {
		return model.NewBadRequestApiError(invalidEmailErrorMessage)
	}
	return nil
}

func validPassword(password string) *model.Error {
	if len(password) < 8 {
		return model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	return nil
}

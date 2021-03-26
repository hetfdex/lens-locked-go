package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"regexp"
	"strings"
)

const pepper = "6Sk65RHhGW7S4qnVPV7m"
const invalidPasswordErrorMessage = "invalid password"

var emailRegex = regexp.MustCompile(`^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$`)

func generateHashFromPassword(password string) (string, *model.Error) {
	if password == "" {
		return "", model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("password"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(hash), nil
}

func compareHashAndPassword(hash string, password string) *model.Error {

	if hash == "" {
		return model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("hash"))
	}

	if password == "" {
		return model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("password"))
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+pepper))

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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}
	return true
}

func validPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"lens-locked-go/util"
	"regexp"
	"strings"
)

const invalidPasswordLengthErrorMessage = "password must be at least 8 characters"

var emailRegex = regexp.MustCompile(`^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$`)

func lower(s string) string {
	return strings.ToLower(s)
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func generateHashFromPassword(password string, pepper string) (string, *model.Error) {
	if password == "" {
		return "", model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}

	hs, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(hs), nil
}

func compareHashAndPassword(hash string, password string, pepper string) *model.Error {

	if hash == "" {
		return model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("hash"))
	}

	if password == "" {
		return model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+pepper))

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(util.InvalidErrorMessage("password"))
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
		return model.NewBadRequestApiError(util.InvalidErrorMessage("email address"))
	}

	if !emailRegex.MatchString(email) {
		return model.NewBadRequestApiError(util.InvalidErrorMessage("email address"))
	}
	return nil
}

func validPassword(password string) *model.Error {
	if len(password) < 8 {
		return model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	return nil
}

package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"lens-locked-go/util"
	"strings"
)

func generateFromPassword(password string) (string, *model.ApiError) {
	if password == "" {
		return "", model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	pw := []byte(password + util.Pepper)

	pwHash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	return string(pwHash), nil
}

func compareHashAndPassword(passwordHash string, password string) *model.ApiError {

	if passwordHash == "" {
		return model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("passwordHash"))
	}

	if password == "" {
		return model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("password"))
	}
	pwHash := []byte(passwordHash)
	pw := []byte(password + util.Pepper)

	err := bcrypt.CompareHashAndPassword(pwHash, pw)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return model.NewForbiddenApiError(util.InvalidPasswordErrorMessage)
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
	return strings.ToLower(strings.TrimSpace(email))
}

func validEmail(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	if !util.EmailRegex.MatchString(email) {
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

package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/rand"
)

func generateFromPassword(user *model.User) *model.ApiError {
	pw := []byte(user.Password + config.Pepper)

	pwHash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	user.Password = ""
	user.PasswordHash = string(pwHash)

	return nil
}

func compareHashAndPassword(user *model.User, password string) *model.ApiError {
	pwHash := []byte(user.PasswordHash)
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

func generateTokenHash(hs *hash.Hasher, token string) (string, *model.ApiError) {
	tokenHash, err := hs.GenerateHash(token)

	if err != nil {
		return "", err
	}
	return tokenHash, nil
}

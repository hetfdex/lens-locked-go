package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/model"
)

func generatePasswordHash(user *model.User) error {
	bytes := []byte(user.Password + pepper)

	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = ""
	user.PasswordHash = string(hash)

	return nil
}

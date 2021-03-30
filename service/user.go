package service

import (
	"golang.org/x/crypto/bcrypt"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/rand"
	"lens-locked-go/repository"
	"regexp"
	"strings"
)

const pepper = "6Sk65RHhGW7S4qnVPV7m"

const key = "yzzmGPkAA9FTmbtzz9jB"

const invalidEmailErrorMessage = "invalid email address"
const invalidPasswordErrorMessage = "invalid password"
const emailInUseErrorMessage = "email address is already in use"
const invalidPasswordLengthErrorMessage = "password must be at least 8 characters"

var emailRegex = regexp.MustCompile(`^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$`)

type IUserService interface {
	Register(*model.RegisterUser) (*model.User, string, *model.Error)
	LoginWithPassword(*model.LoginUser) (*model.User, string, *model.Error)
	LoginWithToken(string) (*model.User, *model.Error)
	Logout(*model.User) *model.Error
}

type userService struct {
	repository repository.IUserRepository
	*hash.Hasher
}

func NewUserService(ur repository.IUserRepository) *userService {
	hs, err := hash.New(key)

	if err != nil {
		panic(err)
	}
	return &userService{
		ur,
		hs,
	}
}

func (s *userService) Register(form *model.RegisterUser) (*model.User, string, *model.Error) {
	form.Email = lower(trimSpace(form.Email))

	if !validEmail(form.Email) {
		return nil, "", model.NewBadRequestApiError(invalidEmailErrorMessage)
	}
	userByEmail, _ := s.getByEmail(form.Email)

	if userByEmail != nil {
		return nil, "", model.NewConflictApiError(emailInUseErrorMessage)
	}

	if !validPassword(form.Password) {
		return nil, "", model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	pwHash, err := generateHashFromPassword(form.Password)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := s.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, "", err
	}
	user := &model.User{
		Name:         form.Name,
		Email:        form.Email,
		PasswordHash: pwHash,
		TokenHash:    tokenHash,
	}

	err = s.repository.Create(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) LoginWithPassword(form *model.LoginUser) (*model.User, string, *model.Error) {
	form.Email = lower(trimSpace(form.Email))

	if !validEmail(form.Email) {
		return nil, "", model.NewBadRequestApiError(invalidEmailErrorMessage)
	}
	user, err := s.getByEmail(form.Email)

	if err != nil {
		return nil, "", err
	}

	if !validPassword(form.Password) {
		return nil, "", model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	err = compareHashAndPassword(user.PasswordHash, form.Password)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := s.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, "", err
	}
	user.TokenHash = tokenHash

	err = s.repository.Update(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) LoginWithToken(token string) (*model.User, *model.Error) {
	tokenHash, err := s.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}
	user, err := s.getByTokenHash(tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Logout(user *model.User) *model.Error {
	user.TokenHash = ""

	err := s.repository.Update(user)

	if err != nil {
		return err
	}
	return nil
}

func (s *userService) getByEmail(email string) (*model.User, *model.Error) {
	if email == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("email"))
	}
	user, err := s.repository.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) getByTokenHash(hash string) (*model.User, *model.Error) {
	if hash == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("tokenHash"))
	}
	user, err := s.repository.Read("token_hash", hash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func generateHashFromPassword(password string) (string, *model.Error) {
	if password == "" {
		return "", model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("password"))
	}

	hs, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)

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

func lower(s string) string {
	return strings.ToLower(s)
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

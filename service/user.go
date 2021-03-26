package service

import (
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const key = "yzzmGPkAA9FTmbtzz9jB"
const invalidEmailErrorMessage = "invalid email address"
const emailInUseErrorMessage = "email address is already in use"
const invalidPasswordLengthErrorMessage = "password must be at least 8 characters"

type IUserService interface {
	Register(*model.UserRegister) (*model.User, string, *model.Error)
	LoginWithPassword(*model.UserLogin) (*model.User, string, *model.Error)
	LoginWithToken(string) (*model.User, *model.Error)
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

func (s *userService) Register(register *model.UserRegister) (*model.User, string, *model.Error) {
	register.Email = lower(trimSpace(register.Email))

	if !validEmail(register.Email) {
		return nil, "", model.NewBadRequestApiError(invalidEmailErrorMessage)
	}
	user, _ := s.getByEmail(register.Email)

	if user != nil {
		return nil, "", model.NewConflictApiError(emailInUseErrorMessage)
	}

	if !validPassword(register.Password) {
		return nil, "", model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	pwHash, err := generateHashFromPassword(register.Password)

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
	user = register.User(pwHash, tokenHash)

	err = s.repository.Create(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) LoginWithPassword(login *model.UserLogin) (*model.User, string, *model.Error) {
	login.Email = lower(trimSpace(login.Email))

	if !validEmail(login.Email) {
		return nil, "", model.NewBadRequestApiError(invalidEmailErrorMessage)
	}
	user, err := s.getByEmail(login.Email)

	if err != nil {
		return nil, "", err
	}

	if !validPassword(login.Password) {
		return nil, "", model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	err = compareHashAndPassword(user.PasswordHash, login.Password)

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

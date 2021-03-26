package service

import (
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const hasherKey = "yzzmGPkAA9FTmbtzz9jB"
const invalidEmailErrorMessage = "invalid email address"
const emailInUseErrorMessage = "email address is already in use"
const invalidPasswordLengthErrorMessage = "password must be at least 8 characters"
const noUserUpdateNeededErrorMessage = "no user update needed"

type IUserService interface {
	Register(*model.RegisterView) (*model.User, string, *model.Error)
	Edit(*model.UpdateView, string) (*model.User, string, *model.Error)
	LoginWithPassword(*model.LoginView) (*model.User, string, *model.Error)
	LoginWithToken(string) (*model.User, *model.Error)
}

type userService struct {
	repository.IUserRepository
	*hash.Hasher
}

func NewUserService(ur repository.IUserRepository) *userService {
	hs, err := hash.New(hasherKey)

	if err != nil {
		panic(err)
	}
	return &userService{
		ur,
		hs,
	}
}

func (s *userService) Register(register *model.RegisterView) (*model.User, string, *model.Error) {
	register.Email = normalizeEmail(register.Email)

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
	pwHash, err := generateFromPassword(register.Password)

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

	err = s.Create(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *userService) Edit(update *model.UpdateView, token string) (*model.User, string, *model.Error) {
	update.Email = normalizeEmail(update.Email)

	if !validEmail(update.Email) {
		return nil, "", model.NewBadRequestApiError(invalidEmailErrorMessage)
	}

	if !validPassword(update.Password) {
		return nil, "", model.NewBadRequestApiError(invalidPasswordLengthErrorMessage)
	}
	user, err := s.LoginWithToken(token)

	if err != nil {
		return nil, "", err
	}
	userFromEmail, _ := s.getByEmail(update.Email)

	if userFromEmail != nil {
		if user.Equals(userFromEmail) {
			return nil, "", model.NewBadRequestApiError(noUserUpdateNeededErrorMessage)
		}
		return nil, "", model.NewConflictApiError(emailInUseErrorMessage)
	}
	newPwHash, err := generateFromPassword(update.Password)

	if err != nil {
		return nil, "", err
	}
	newToken, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	newTokenHash, err := s.Hasher.GenerateTokenHash(newToken)

	if err != nil {
		return nil, "", err
	}
	user.Update(update, newPwHash, newTokenHash)

	err = s.Update(user)

	if err != nil {
		return nil, "", err
	}
	return user, newToken, nil
}

func (s *userService) LoginWithPassword(login *model.LoginView) (*model.User, string, *model.Error) {
	login.Email = normalizeEmail(login.Email)

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

	err = s.Update(user)

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
	user, err := s.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) getByTokenHash(tokenHash string) (*model.User, *model.Error) {
	if tokenHash == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("tokenHash"))
	}
	user, err := s.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

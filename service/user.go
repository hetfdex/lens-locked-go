package service

import (
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

type IUserService interface {
	LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError)
	LoginWithToken(token string) (*model.User, *model.ApiError)
	Register(user *model.User) *model.ApiError
	GetByEmail(email string) (*model.User, *model.ApiError)
	GetByTokenHash(tokenHash string) (*model.User, *model.ApiError)
	UpdateToken(user *model.User) *model.ApiError
}

type UserService struct {
	repository.IUserRepository
	*hash.Hasher
}

func NewUserService(ur repository.IUserRepository) *UserService {
	hs := hash.New(config.HasherKey)

	return &UserService{
		ur,
		hs,
	}
}

func (us *UserService) LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError) {
	user, err := us.GetByEmail(login.Email)

	if err != nil {
		return nil, err
	}
	err = compareHashAndPassword(user, login.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) LoginWithToken(token string) (*model.User, *model.ApiError) {
	tokenHash, err := generateTokenHash(us.Hasher, token)

	if err != nil {
		return nil, err
	}
	user, err := us.GetByTokenHash(tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) Register(user *model.User) *model.ApiError {
	existingUser, _ := us.GetByEmail(user.Email)

	if existingUser != nil {
		return model.NewConflictApiError("user already exists")
	}
	err := generateFromPassword(user)

	if err != nil {
		return err
	}
	token, err := generateToken()

	if err != nil {
		return err
	}
	user.Token = token
	tokenHash, err := generateTokenHash(us.Hasher, user.Token)

	if err != nil {
		return err
	}
	user.TokenHash = tokenHash

	err = us.Create(user)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetByEmail(email string) (*model.User, *model.ApiError) {
	if email == "" {
		return nil, model.NewInternalServerApiError("string must not be empty")
	}
	user, err := us.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetByTokenHash(tokenHash string) (*model.User, *model.ApiError) {
	if tokenHash == "" {
		return nil, model.NewInternalServerApiError("string must not be empty")
	}
	user, err := us.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdateToken(user *model.User) *model.ApiError {
	token, err := generateToken()

	if err != nil {
		return err
	}
	user.Token = token

	tokenHash, err := generateTokenHash(us.Hasher, user.Token)

	if err != nil {
		return err
	}
	user.TokenHash = tokenHash

	err = us.Update(user)

	if err != nil {
		return err
	}
	return nil
}

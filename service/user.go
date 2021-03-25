package service

import (
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/validator"
)

type IUserService interface {
	Register(register *model.RegisterForm) (*model.User, *model.ApiError)
	LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError)
	LoginWithToken(token string) (*model.User, *model.ApiError)
	GetByEmail(email string) (*model.User, *model.ApiError)
	GetByTokenHash(tokenHash string) (*model.User, *model.ApiError)
}

type userService struct {
	repository.IUserRepository
	*hash.Hasher
}

func NewUserService(ur repository.IUserRepository) *userService {
	hs, err := hash.New(config.HasherKey)

	if err != nil {
		panic(err)
	}
	return &userService{
		ur,
		hs,
	}
}

func (us *userService) Register(register *model.RegisterForm) (*model.User, *model.ApiError) {
	user, _ := us.GetByEmail(register.Email)

	if user != nil {
		return nil, model.NewConflictApiError("user already exists")
	}
	pwHash, err := generateFromPassword(register.Password)

	if err != nil {
		return nil, err
	}
	user = &model.User{
		Name:         register.Name,
		Email:        register.Email,
		Password:     "",
		PasswordHash: pwHash,
		Token:        "",
		TokenHash:    "",
	}

	token, err := generateToken()

	if err != nil {
		return nil, err
	}
	user.Token = token

	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}
	user.TokenHash = tokenHash

	err = us.Create(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError) {
	user, err := us.GetByEmail(login.Email)

	if err != nil {
		return nil, err
	}
	err = compareHashAndPassword(user.PasswordHash, login.Password)

	if err != nil {
		return nil, err
	}
	token, err := generateToken()

	if err != nil {
		return nil, err
	}
	user.Token = token

	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}
	user.TokenHash = tokenHash

	err = us.Update(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) LoginWithToken(token string) (*model.User, *model.ApiError) {
	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}
	user, err := us.GetByTokenHash(tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetByEmail(email string) (*model.User, *model.ApiError) {
	if validator.EmptyString(email) {
		return nil, model.NewInternalServerApiError("email must not be empty")
	}
	user, err := us.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetByTokenHash(tokenHash string) (*model.User, *model.ApiError) {
	if validator.EmptyString(tokenHash) {
		return nil, model.NewInternalServerApiError("tokenHash must not be empty")
	}
	user, err := us.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

package service

import (
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

type IUserService interface {
	Register(register *model.RegisterForm) (*model.User, *model.ApiError)
	LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError)
	LoginWithToken(token string) (*model.User, *model.ApiError)
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
	register.Email = normalizeEmail(register.Email)

	user, _ := us.getByEmail(register.Email)

	if user != nil {
		return nil, model.NewConflictApiError("user already exists")
	}
	pwHash, err := generateFromPassword(register.Password)

	if err != nil {
		return nil, err
	}

	token, err := generateToken()

	if err != nil {
		return nil, err
	}

	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}

	user = model.NewUserFromRegister(register, pwHash, token, tokenHash)

	err = us.Create(user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError) {
	login.Email = normalizeEmail(login.Email)

	user, err := us.getByEmail(login.Email)

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
	user, err := us.getByTokenHash(tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) getByEmail(email string) (*model.User, *model.ApiError) {
	if email == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("email"))
	}
	user, err := us.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) getByTokenHash(tokenHash string) (*model.User, *model.ApiError) {
	if tokenHash == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("tokenHash"))
	}
	user, err := us.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

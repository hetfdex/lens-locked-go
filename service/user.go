package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
	hasher         *hash.Hasher
}

func NewUserService(ur *repository.UserRepository, hs *hash.Hasher) (*UserService, *model.ApiError) {
	return &UserService{
		userRepository: ur,
		hasher:         hs,
	}, nil
}

func (us *UserService) LoginWithPassword(login *model.LoginForm) (*model.User, *model.ApiError) {
	user, err := us.GetUserByEmail(login.Email)

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
	tokenHash, err := generateTokenHash(us.hasher, token)

	if err != nil {
		return nil, err
	}
	user, err := us.GetUserByTokenHash(tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) RegisterUser(user *model.User) *model.ApiError {
	existingUser, _ := us.GetUserByEmail(user.Email)

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
	tokenHash, err := generateTokenHash(us.hasher, user.Token)

	if err != nil {
		return err
	}
	user.TokenHash = tokenHash

	err = us.userRepository.Create(user)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserByEmail(email string) (*model.User, *model.ApiError) {
	if email == "" {
		return nil, model.NewInternalServerApiError("string must not be empty")
	}
	user, err := us.userRepository.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUserByTokenHash(tokenHash string) (*model.User, *model.ApiError) {
	if tokenHash == "" {
		return nil, model.NewInternalServerApiError("string must not be empty")
	}
	user, err := us.userRepository.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) UpdateUserToken(user *model.User) *model.ApiError {
	token, err := generateToken()

	if err != nil {
		return err
	}
	user.Token = token

	tokenHash, err := generateTokenHash(us.hasher, user.Token)

	if err != nil {
		return err
	}
	user.TokenHash = tokenHash

	err = us.UpdateUser(user)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateUser(user *model.User) *model.ApiError {
	err := us.userRepository.Update(user)

	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) Delete(id uuid.UUID) *model.ApiError {
	err := us.userRepository.Delete(id)

	if err != nil {
		return err
	}
	return nil
}

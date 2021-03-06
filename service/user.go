package service

import (
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
)

type IUserService interface {
	Register(*model.RegisterUser) (*model.User, string, *model.Error)
	LoginWithPassword(*model.LoginUser) (*model.User, string, *model.Error)
	LoginWithToken(string) (*model.User, *model.Error)
	Logout(*model.User) *model.Error
}

type userService struct {
	repository   repository.IUserRepository
	cryptoConfig *config.CryptoConfig
	hasher       *hash.Hasher
}

func newUserService(ur repository.IUserRepository, cc *config.CryptoConfig) *userService {
	hs, err := hash.New(cc.HasherKey)

	if err != nil {
		panic(err)
	}
	return &userService{
		repository:   ur,
		cryptoConfig: cc,
		hasher:       hs,
	}
}

func (s *userService) Register(form *model.RegisterUser) (*model.User, string, *model.Error) {
	form.Email = lower(trimSpace(form.Email))

	err := validEmail(form.Email)

	if err != nil {
		return nil, "", err
	}
	userByEmail, _ := s.getByEmail(form.Email)

	if userByEmail != nil {
		return nil, "", model.NewConflictApiError(util.InUseErrorMessage("email address"))
	}
	err = validPassword(form.Password)

	if err != nil {
		return nil, "", err
	}
	pwHash, err := generateHashFromPassword(form.Password, s.cryptoConfig.Pepper)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := s.hasher.GenerateTokenHash(token)

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

	err := validEmail(form.Email)

	if err != nil {
		return nil, "", err
	}
	user, err := s.getByEmail(form.Email)

	if err != nil {
		return nil, "", err
	}

	err = validPassword(form.Password)

	if err != nil {
		return nil, "", err
	}
	err = compareHashAndPassword(user.PasswordHash, form.Password, s.cryptoConfig.Pepper)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := s.hasher.GenerateTokenHash(token)

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
	tokenHash, err := s.hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, err
	}
	user, err := s.repository.Read("token_hash", tokenHash)

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
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}
	user, err := s.repository.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

package service

import (
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
)

type IUserService interface {
	Register(register *model.Register) (*model.User, string, *model.ApiError)
	Edit(update *model.Update, token string) (*model.User, string, *model.ApiError)
	LoginWithPassword(login *model.Login) (*model.User, string, *model.ApiError)
	LoginWithToken(token string) (*model.User, *model.ApiError)
}

type userService struct {
	repository.IUserRepository
	*hash.Hasher
}

func NewUserService(ur repository.IUserRepository) *userService {
	hs, err := hash.New(util.HasherKey)

	if err != nil {
		panic(err)
	}
	return &userService{
		ur,
		hs,
	}
}

func (us *userService) Register(register *model.Register) (*model.User, string, *model.ApiError) {
	register.Email = normalizeEmail(register.Email)

	if !validEmail(register.Email) {
		return nil, "", model.NewBadRequestApiError(util.InvalidEmailErrorMessage)
	}
	user, _ := us.getByEmail(register.Email)

	if user != nil {
		return nil, "", model.NewConflictApiError(util.EmailInUseErrorMessage)
	}

	if !validPassword(register.Password) {
		return nil, "", model.NewBadRequestApiError(util.InvalidPasswordErrorMessage)
	}
	pwHash, err := generateFromPassword(register.Password)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, "", err
	}
	user = register.User(pwHash, tokenHash)

	err = us.Create(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (us *userService) Edit(update *model.Update, token string) (*model.User, string, *model.ApiError) {
	update.Email = normalizeEmail(update.Email)

	if !validEmail(update.Email) {
		return nil, "", model.NewBadRequestApiError(util.InvalidEmailErrorMessage)
	}

	if !validPassword(update.Password) {
		return nil, "", model.NewBadRequestApiError(util.InvalidPasswordErrorMessage)
	}
	user, err := us.LoginWithToken(token)

	if err != nil {
		return nil, "", err
	}
	userFromEmail, _ := us.getByEmail(update.Email)

	if userFromEmail != nil {
		if user.Equals(userFromEmail) {
			return nil, "", model.NewBadRequestApiError(util.NoUserUpdateNeededErrorMessage)
		}
		return nil, "", model.NewConflictApiError(util.EmailInUseErrorMessage)
	}
	newPwHash, err := generateFromPassword(update.Password)

	if err != nil {
		return nil, "", err
	}
	newToken, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	newTokenHash, err := us.Hasher.GenerateTokenHash(newToken)

	if err != nil {
		return nil, "", err
	}
	user.Update(update, newPwHash, newTokenHash)

	err = us.Update(user)

	if err != nil {
		return nil, "", err
	}
	return user, newToken, nil
}

func (us *userService) LoginWithPassword(login *model.Login) (*model.User, string, *model.ApiError) {
	login.Email = normalizeEmail(login.Email)

	if !validEmail(login.Email) {
		return nil, "", model.NewBadRequestApiError(util.InvalidEmailErrorMessage)
	}
	user, err := us.getByEmail(login.Email)

	if err != nil {
		return nil, "", err
	}

	if !validPassword(login.Password) {
		return nil, "", model.NewBadRequestApiError(util.InvalidPasswordErrorMessage)
	}
	err = compareHashAndPassword(user.PasswordHash, login.Password)

	if err != nil {
		return nil, "", err
	}
	token, err := generateToken()

	if err != nil {
		return nil, "", err
	}
	tokenHash, err := us.Hasher.GenerateTokenHash(token)

	if err != nil {
		return nil, "", err
	}
	user.TokenHash = tokenHash

	err = us.Update(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
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
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("email"))
	}
	user, err := us.Read("email", email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) getByTokenHash(tokenHash string) (*model.User, *model.ApiError) {
	if tokenHash == "" {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("tokenHash"))
	}
	user, err := us.Read("token_hash", tokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

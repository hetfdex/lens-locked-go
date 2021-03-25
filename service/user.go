package service

import (
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

type IUserService interface {
	Register(register *model.RegisterForm) (*model.User, string, *model.ApiError)
	Edit(update *model.UpdateForm, token string) (*model.User, string, *model.ApiError)
	LoginWithPassword(login *model.LoginForm) (*model.User, string, *model.ApiError)
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

func (us *userService) Register(register *model.RegisterForm) (*model.User, string, *model.ApiError) {
	register.Email = normalizeEmail(register.Email)

	if !validEmail(register.Email) {
		return nil, "", model.NewBadRequestApiError("invalid email")
	}
	user, _ := us.getByEmail(register.Email)

	if user != nil {
		return nil, "", model.NewConflictApiError("email is already registered")
	}

	if !validPassword(register.Password) {
		return nil, "", model.NewBadRequestApiError("invalid password")
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
	user = model.NewUserFromRegister(register, pwHash, tokenHash)

	err = us.Create(user)

	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (us *userService) Edit(update *model.UpdateForm, token string) (*model.User, string, *model.ApiError) {
	update.Email = normalizeEmail(update.Email)

	if !validEmail(update.Email) {
		return nil, "", model.NewBadRequestApiError("invalid email")
	}

	if !validPassword(update.Password) {
		return nil, "", model.NewBadRequestApiError("invalid password")
	}
	loggedInUser, err := us.LoginWithToken(token)

	if err != nil {
		return nil, "", err
	}
	userFromEmail, _ := us.getByEmail(update.Email)

	if userFromEmail != nil {
		if loggedInUser.DeepEqual(userFromEmail) {
			return nil, "", model.NewBadRequestApiError("no user data changed")
		}
		return nil, "", model.NewConflictApiError("email is already registered")
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
	updatedUser := model.NewUserFromUpdate(update, loggedInUser, newPwHash, newTokenHash)

	err = us.Update(updatedUser)

	if err != nil {
		return nil, "", err
	}
	return updatedUser, newToken, nil
}

func (us *userService) LoginWithPassword(login *model.LoginForm) (*model.User, string, *model.ApiError) {
	login.Email = normalizeEmail(login.Email)

	if !validEmail(login.Email) {
		return nil, "", model.NewBadRequestApiError("invalid email")
	}
	user, err := us.getByEmail(login.Email)

	if err != nil {
		return nil, "", err
	}

	if !validPassword(login.Password) {
		return nil, "", model.NewBadRequestApiError("invalid password")
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

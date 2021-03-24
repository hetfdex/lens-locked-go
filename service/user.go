package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/config"
	"lens-locked-go/hash"
	"lens-locked-go/model"
)

type UserService struct {
	db     *gorm.DB
	hasher *hash.Hasher
}

func New(dsn string) (*UserService, *model.ApiError) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	hasher := hash.New(config.HasherKey)

	return &UserService{
		db:     db,
		hasher: hasher,
	}, nil
}

func (u *UserService) CreateTable() *model.ApiError {
	err := u.db.Migrator().CreateTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (u *UserService) DropTable() *model.ApiError {
	err := u.db.Migrator().DropTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (u *UserService) AuthenticateWithPassword(login *model.LoginForm) (*model.User, *model.ApiError) {
	user, err := u.Read("email", login.Email)

	if err != nil {
		return nil, err
	}
	err = compareHashAndPassword(user, login.Password)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) AuthenticateWithToken(token string) (*model.User, *model.ApiError) {
	user := &model.User{}

	user.Token = token

	err := generateTokenHash(u.hasher, user)

	if err != nil {
		return nil, err
	}
	user, err = u.Read("token_hash", user.TokenHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) Create(user *model.User) *model.ApiError {
	apiErr := generateFromPassword(user)

	if user.Token == "" {
		apiErr = generateToken(user)

		if apiErr != nil {
			return apiErr
		}
	}
	apiErr = generateTokenHash(u.hasher, user)

	if apiErr != nil {
		return apiErr
	}
	err := u.db.Create(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (u *UserService) Read(field string, value interface{}) (*model.User, *model.ApiError) {
	if field == "" {
		return nil, model.NewInternalServerApiError("string must not be empty")
	}
	user := &model.User{}

	cond := fmt.Sprintf("%s = ?", field)

	err := u.db.First(user, cond, value).Error

	if err != nil {
		return nil, model.NewNotFoundApiError(err.Error())
	}

	return user, nil
}

func (u *UserService) Update(user *model.User) *model.ApiError {
	err := u.db.Save(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (u *UserService) Delete(id uuid.UUID) *model.ApiError {
	err := u.db.Delete(&model.User{}, id).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

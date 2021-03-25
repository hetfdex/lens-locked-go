package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
	"lens-locked-go/util"
)

type IUserRepository interface {
	Create(user *model.User) *model.ApiError
	Read(field string, value interface{}) (*model.User, *model.ApiError)
	Update(user *model.User) *model.ApiError
	Delete(user *model.User) *model.ApiError
}

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		database: db,
	}
}

func (ur *userRepository) Create(user *model.User) *model.ApiError {
	err := ur.database.Create(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *userRepository) Read(field string, value interface{}) (*model.User, *model.ApiError) {

	if field == "" {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("field"))
	}
	user := &model.User{}

	query := fmt.Sprintf("%s = ?", field)

	err := ur.database.First(user, query, value).Error

	if err != nil {
		return nil, model.NewNotFoundApiError(err.Error())
	}
	return user, nil
}

func (ur *userRepository) Update(user *model.User) *model.ApiError {
	err := ur.database.Save(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *userRepository) Delete(user *model.User) *model.ApiError {
	err := ur.database.Delete(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

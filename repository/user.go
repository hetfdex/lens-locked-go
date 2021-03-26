package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

const userNotFoundError = "user not found"

type IUserRepository interface {
	Create(*model.User) *model.Error
	Read(string, interface{}) (*model.User, *model.Error)
	Update(*model.User) *model.Error
	Delete(*model.User) *model.Error
}

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		database: db,
	}
}

func (r *userRepository) Create(user *model.User) *model.Error {
	err := r.database.Create(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *userRepository) Read(field string, value interface{}) (*model.User, *model.Error) {

	if field == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("field"))
	}
	user := &model.User{}

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.First(user, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(userNotFoundError)
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return user, nil
}

func (r *userRepository) Update(user *model.User) *model.Error {
	err := r.database.Save(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *userRepository) Delete(user *model.User) *model.Error {
	err := r.database.Delete(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

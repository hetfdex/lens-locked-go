package repository

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) (*UserRepository, *model.ApiError) {
	return &UserRepository{
		database: db,
	}, nil
}

func (ur *UserRepository) CreateTable() *model.ApiError {
	err := ur.database.Migrator().CreateTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) DropTable() *model.ApiError {
	err := ur.database.Migrator().DropTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Create(user *model.User) *model.ApiError {
	err := ur.database.Create(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Read(field string, value interface{}) (*model.User, *model.ApiError) {
	user := &model.User{}

	query := fmt.Sprintf("%s = ?", field)

	err := ur.database.First(user, query, value).Error

	if err != nil {
		return nil, model.NewNotFoundApiError(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) Update(user *model.User) *model.ApiError {
	err := ur.database.Save(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Delete(id uuid.UUID) *model.ApiError {
	err := ur.database.Delete(&model.User{}, id).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

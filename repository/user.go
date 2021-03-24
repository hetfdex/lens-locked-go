package repository

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dsn string) (*UserRepository, *model.ApiError) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return &UserRepository{
		db: db,
	}, nil
}

func (ur *UserRepository) CreateTable() *model.ApiError {
	err := ur.db.Migrator().CreateTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) DropTable() *model.ApiError {
	err := ur.db.Migrator().DropTable(&model.User{})

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Create(user *model.User) *model.ApiError {
	err := ur.db.Create(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Read(field string, value interface{}) (*model.User, *model.ApiError) {
	user := &model.User{}

	cond := fmt.Sprintf("%s = ?", field)

	err := ur.db.First(user, cond, value).Error

	if err != nil {
		return nil, model.NewNotFoundApiError(err.Error())
	}
	return user, nil
}

func (ur *UserRepository) Update(user *model.User) *model.ApiError {
	err := ur.db.Save(user).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (ur *UserRepository) Delete(id uuid.UUID) *model.ApiError {
	err := ur.db.Delete(&model.User{}, id).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

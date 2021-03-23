package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

type UserService struct {
	db *gorm.DB
}

func New(dsn string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

func (u *UserService) CreateTable() error {
	err := u.db.Migrator().CreateTable(&model.User{})

	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DropTable() error {
	err := u.db.Migrator().DropTable(&model.User{})

	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) Create(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *UserService) Read(field string, value interface{}) (*model.User, error) {
	user := &model.User{}

	cond := fmt.Sprintf("%s = ?", field)

	err := u.db.First(user, cond, value).Error

	return user, err
}

func (u *UserService) Update(user *model.User) error {
	return u.db.Save(user).Error
}

func (u *UserService) Delete(id uuid.UUID) error {
	return u.db.Delete(&model.User{}, id).Error
}

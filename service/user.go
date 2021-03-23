package service

import (
	"github.com/gofrs/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

type userService struct {
	db *gorm.DB
}

func New(dsn string) (*userService, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return &userService{
		db: db,
	}, nil
}

func (u *userService) CreateTable() error {
	err := u.db.Migrator().CreateTable(&model.User{})

	if err != nil {
		return err
	}
	return nil
}

func (u *userService) DropTable() error {
	err := u.db.Migrator().DropTable(&model.User{})

	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Create(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *userService) ById(id uuid.UUID) (*model.User, error) {
	var user model.User

	err := u.db.First(&user, "id = ?", id).Error

	return &user, err
}

func (u *userService) ByEmail(email string) (*model.User, error) {
	var user model.User

	err := u.db.First(&user, "email = ?", email).Error

	return &user, err
}

func (u *userService) Update(user *model.User) error {
	return u.db.Save(user).Error
}

func (u *userService) Delete(id uuid.UUID) error {
	return u.db.Delete(&model.User{}, id).Error
}

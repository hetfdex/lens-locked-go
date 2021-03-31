package repository

import (
	"gorm.io/gorm"
	"lens-locked-go/model"
)

type IImageRepository interface {
	Create(*model.Image) *model.Error
}

type imageRepository struct {
	database *gorm.DB
}

func NewImageRepository(db *gorm.DB) *imageRepository {
	return &imageRepository{
		database: db,
	}
}

func (r *imageRepository) Create(image *model.Image) *model.Error {
	err := r.database.Create(image).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

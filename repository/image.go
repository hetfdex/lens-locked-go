package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

const imageNotFoundError = "image not found"

type IImageRepository interface {
	Create(*model.Image) *model.Error
	Read(string, interface{}) (*model.Image, *model.Error)
	ReadAll(string, interface{}) ([]*model.Image, *model.Error)
	Update(*model.Image) *model.Error
	Delete(*model.Image) *model.Error
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

func (r *imageRepository) Read(field string, value interface{}) (*model.Image, *model.Error) {
	if field == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("field"))
	}
	image := &model.Image{}

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.First(image, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(imageNotFoundError)
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return image, nil
}

func (r *imageRepository) ReadAll(field string, value interface{}) ([]*model.Image, *model.Error) {
	if field == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("field"))
	}
	var images []*model.Image

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.Find(&images, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(imageNotFoundError)
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return images, nil
}

func (r *imageRepository) Update(image *model.Image) *model.Error {
	err := r.database.Save(image).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *imageRepository) Delete(image *model.Image) *model.Error {
	err := r.database.Delete(image).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

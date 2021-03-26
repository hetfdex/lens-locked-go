package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

const galleryNotFoundError = "gallery not found"

type IGalleryRepository interface {
	Create(gallery *model.Gallery) *model.ApiError
	Read(field string, value interface{}) (*model.Gallery, *model.ApiError)
	Update(gallery *model.Gallery) *model.ApiError
	Delete(gallery *model.Gallery) *model.ApiError
}

type galleryRepository struct {
	database *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *galleryRepository {
	return &galleryRepository{
		database: db,
	}
}

func (gr *galleryRepository) Create(gallery *model.Gallery) *model.ApiError {
	err := gr.database.Create(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (gr *galleryRepository) Read(field string, value interface{}) (*model.Gallery, *model.ApiError) {

	if field == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("field"))
	}
	gallery := &model.Gallery{}

	query := fmt.Sprintf("%s = ?", field)

	err := gr.database.First(gallery, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(galleryNotFoundError)
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return gallery, nil
}

func (gr *galleryRepository) Update(gallery *model.Gallery) *model.ApiError {
	err := gr.database.Save(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (gr *galleryRepository) Delete(gallery *model.Gallery) *model.ApiError {
	err := gr.database.Delete(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
)

const galleryNotFoundError = "gallery not found"

type IGalleryRepository interface {
	Create(*model.Gallery) *model.Error
	Read(string, interface{}) (*model.Gallery, *model.Error)
	Update(*model.Gallery) *model.Error
	Delete(*model.Gallery) *model.Error
}

type galleryRepository struct {
	database *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *galleryRepository {
	return &galleryRepository{
		database: db,
	}
}

func (r *galleryRepository) Create(gallery *model.Gallery) *model.Error {
	err := r.database.Create(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *galleryRepository) Read(field string, value interface{}) (*model.Gallery, *model.Error) {

	if field == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("field"))
	}
	gallery := &model.Gallery{}

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.First(gallery, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(galleryNotFoundError)
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return gallery, nil
}

func (r *galleryRepository) Update(gallery *model.Gallery) *model.Error {
	err := r.database.Save(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *galleryRepository) Delete(gallery *model.Gallery) *model.Error {
	err := r.database.Delete(gallery).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

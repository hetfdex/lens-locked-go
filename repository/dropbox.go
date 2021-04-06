package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
	"lens-locked-go/util"
)

type IDropboxRepository interface {
	Create(*model.Dropbox) *model.Error
	Read(string, interface{}) (*model.Dropbox, *model.Error)
	Delete(*model.Dropbox) *model.Error
}

type dropboxRepository struct {
	database *gorm.DB
}

func newDropboxRepository(db *gorm.DB) *dropboxRepository {
	return &dropboxRepository{
		database: db,
	}
}

func (r *dropboxRepository) Create(dropbox *model.Dropbox) *model.Error {
	err := r.database.Create(dropbox).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *dropboxRepository) Read(field string, value interface{}) (*model.Dropbox, *model.Error) {
	if field == "" {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("field"))
	}
	dropbox := &model.Dropbox{}

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.First(dropbox, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(notFoundErrorMessage("dropbox"))
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return dropbox, nil
}

func (r *dropboxRepository) Delete(dropbox *model.Dropbox) *model.Error {
	err := r.database.Unscoped().Delete(dropbox).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

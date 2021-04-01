package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lens-locked-go/model"
	"lens-locked-go/util"
)

type IOAuthRepository interface {
	Create(*model.OAuth) *model.Error
	Read(string, interface{}) (*model.OAuth, *model.Error)
	Delete(*model.OAuth) *model.Error
}

type oAuthRepository struct {
	database *gorm.DB
}

func newOAuthRepository(db *gorm.DB) *oAuthRepository {
	return &oAuthRepository{
		database: db,
	}
}

func (r *oAuthRepository) Create(oauth *model.OAuth) *model.Error {
	err := r.database.Create(oauth).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

func (r *oAuthRepository) Read(field string, value interface{}) (*model.OAuth, *model.Error) {
	if field == "" {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("field"))
	}
	oauth := &model.OAuth{}

	query := fmt.Sprintf("%s = ?", field)

	err := r.database.First(oauth, query, value).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewNotFoundApiError(notFoundErrorMessage("oauth"))
		}
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return oauth, nil
}

func (r *oAuthRepository) Delete(oauth *model.OAuth) *model.Error {
	err := r.database.Delete(oauth).Error

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	return nil
}

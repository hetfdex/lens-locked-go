package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
)

type IDropboxService interface {
	Create(*model.Dropbox) *model.Error
	GetByUserId(uuid.UUID) (*model.Dropbox, *model.Error)
	Delete(*model.Dropbox) *model.Error
}

type dropboxService struct {
	repository repository.IDropboxRepository
}

func newDropboxService(dr repository.IDropboxRepository) *dropboxService {
	return &dropboxService{
		dr,
	}
}

func (s *dropboxService) Create(dropbox *model.Dropbox) *model.Error {
	return s.repository.Create(dropbox)
}

func (s *dropboxService) GetByUserId(userId uuid.UUID) (*model.Dropbox, *model.Error) {
	if userId == uuid.Nil {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("userId"))
	}
	dropbox, err := s.repository.Read("user_id", userId)

	if err != nil {
		return nil, err
	}
	return dropbox, nil
}

func (s *dropboxService) Delete(dropbox *model.Dropbox) *model.Error {
	return s.repository.Delete(dropbox)
}

package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
)

type IOAuthService interface {
	Create(*model.OAuth) *model.Error
	GetByUserId(uuid.UUID) (*model.OAuth, *model.Error)
	Delete(*model.OAuth) *model.Error
}

type oAuthService struct {
	repository repository.IOAuthRepository
}

func newOAuthService(or repository.IOAuthRepository) *oAuthService {
	return &oAuthService{
		or,
	}
}

func (s *oAuthService) Create(oauth *model.OAuth) *model.Error {
	return s.repository.Create(oauth)
}

func (s *oAuthService) GetByUserId(userId uuid.UUID) (*model.OAuth, *model.Error) {
	if userId == uuid.Nil {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("userId"))
	}
	oauth, err := s.repository.Read("user_id", userId)

	if err != nil {
		return nil, err
	}
	return oauth, nil
}

func (s *oAuthService) Delete(oauth *model.OAuth) *model.Error {
	return s.repository.Delete(oauth)
}

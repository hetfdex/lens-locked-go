package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
)

type IGalleryService interface {
	Create(*model.CreateGallery, uuid.UUID) (*model.Gallery, *model.Error)
	GetById(uuid.UUID) (*model.Gallery, *model.Error)
	GetAllByUserId(uuid.UUID) ([]*model.Gallery, *model.Error)
	Update(*model.Gallery, *model.EditGallery) *model.Error
	Delete(*model.Gallery) *model.Error
}

type galleryService struct {
	repository repository.IGalleryRepository
}

func NewGalleryService(ur repository.IGalleryRepository) *galleryService {
	return &galleryService{
		ur,
	}
}

func (s *galleryService) Create(form *model.CreateGallery, userId uuid.UUID) (*model.Gallery, *model.Error) {
	form.Name = trimSpace(form.Name)

	galleriesByName, err := s.getAllByName(form.Name)

	if err != nil {
		return nil, err
	}

	for _, g := range galleriesByName {
		if g.UserId == userId {
			return nil, model.NewConflictApiError(util.InUseErrorMessage("name"))
		}
	}

	gallery := &model.Gallery{
		Name:   form.Name,
		UserId: userId,
	}

	err = s.repository.Create(gallery)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) GetById(id uuid.UUID) (*model.Gallery, *model.Error) {
	if id == uuid.Nil {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("id"))
	}
	gallery, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) GetAllByUserId(userId uuid.UUID) ([]*model.Gallery, *model.Error) {
	if userId == uuid.Nil {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("userId"))
	}
	galleries, err := s.repository.ReadAll("user_id", userId)

	if err != nil {
		return nil, err
	}
	return galleries, nil
}

func (s *galleryService) Update(gallery *model.Gallery, form *model.EditGallery) *model.Error {
	form.Name = trimSpace(form.Name)

	galleriesByName, err := s.getAllByName(form.Name)

	if err != nil {
		return err
	}

	for _, g := range galleriesByName {
		if g.UserId == gallery.UserId {
			return model.NewConflictApiError(util.InUseErrorMessage("name"))
		}
	}

	gallery.Name = form.Name

	err = s.repository.Update(gallery)

	if err != nil {
		return err
	}
	return nil
}

func (s *galleryService) Delete(gallery *model.Gallery) *model.Error {
	return s.repository.Delete(gallery)
}

func (s *galleryService) getAllByName(name string) ([]*model.Gallery, *model.Error) {
	if name == "" {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("name"))
	}
	gallery, err := s.repository.ReadAll("name", name)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

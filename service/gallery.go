package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const titleInUseErrorMessage = "title is already in use"

type IGalleryService interface {
	Create(*model.CreateGallery, uuid.UUID) (*model.Gallery, *model.Error)
	GetById(uuid.UUID) (*model.Gallery, *model.Error)
	GetAllByUserId(uuid.UUID) ([]model.Gallery, *model.Error)
	Edit(*model.Gallery, *model.EditGallery) (*model.Gallery, *model.Error)
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

func (s *galleryService) Create(create *model.CreateGallery, userId uuid.UUID) (*model.Gallery, *model.Error) {
	create.Name = trimSpace(create.Name)

	galleriesByTitle, err := s.getAllByTitle(create.Name)

	if err != nil {
		return nil, err
	}

	if galleriesByTitle != nil {
		for _, g := range galleriesByTitle {
			if g.UserId == userId {
				return nil, model.NewConflictApiError(titleInUseErrorMessage)
			}
		}
	}
	gallery := create.Gallery(userId)

	err = s.repository.Create(gallery)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) GetById(id uuid.UUID) (*model.Gallery, *model.Error) {
	if id == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("id"))
	}
	gallery, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) GetAllByUserId(userId uuid.UUID) ([]model.Gallery, *model.Error) {
	if userId == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("userId"))
	}
	gallery, err := s.repository.ReadAll("user_id", userId)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) Edit(gallery *model.Gallery, edit *model.EditGallery) (*model.Gallery, *model.Error) {
	edit.Name = trimSpace(edit.Name)

	galleriesByTitle, err := s.getAllByTitle(edit.Name)

	if err != nil {
		return nil, err
	}

	if galleriesByTitle != nil {
		for _, g := range galleriesByTitle {
			if g.UserId == gallery.UserId {
				return nil, model.NewConflictApiError(titleInUseErrorMessage)
			}
		}
	}
	gallery = edit.Gallery(gallery.UserId)

	err = s.repository.Update(gallery)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) Delete(gallery *model.Gallery) *model.Error {
	return s.repository.Delete(gallery)
}

func (s *galleryService) getAllByTitle(title string) ([]model.Gallery, *model.Error) {
	if title == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("title"))
	}
	gallery, err := s.repository.ReadAll("title", title)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

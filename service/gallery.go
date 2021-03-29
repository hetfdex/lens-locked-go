package service

import (
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const titleInUseErrorMessage = "title is already in use"

type IGalleryService interface {
	Create(gallery *model.CreateGallery, userId uuid.UUID) (*model.Gallery, *model.Error)
	Get(id uuid.UUID) (*model.Gallery, *model.Error)
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

	gallery, _ := s.getByTitle(create.Name)

	if gallery != nil && gallery.UserId == userId {
		return nil, model.NewConflictApiError(titleInUseErrorMessage)
	}
	gallery = create.Gallery(userId)

	err := s.repository.Create(gallery)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) Get(id uuid.UUID) (*model.Gallery, *model.Error) {
	return s.getById(id)
}

func (s *galleryService) getById(id uuid.UUID) (*model.Gallery, *model.Error) {
	if id == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("id"))
	}
	gallery, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) getByTitle(title string) (*model.Gallery, *model.Error) {
	if title == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("title"))
	}
	gallery, err := s.repository.Read("title", title)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

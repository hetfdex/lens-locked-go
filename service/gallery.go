package service

import (
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const titleInUseErrorMessage = "title is already in use"

type IGalleryService interface {
	Create(gallery *model.CreateGallery) (*model.Gallery, *model.Error)
}

type galleryService struct {
	repository repository.IGalleryRepository
}

func NewGalleryService(ur repository.IGalleryRepository) *galleryService {
	return &galleryService{
		ur,
	}
}

func (s *galleryService) Create(create *model.CreateGallery) (*model.Gallery, *model.Error) {
	create.Name = normalizeEmail(create.Name)

	gallery, _ := s.getByTitle(create.Name)

	if gallery != nil && gallery.UserId.String() == "TODO" {
		return nil, model.NewConflictApiError(titleInUseErrorMessage)
	}
	gallery = create.Gallery()

	err := s.repository.Create(gallery)

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

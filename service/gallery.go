package service

import (
	"lens-locked-go/model"
	"lens-locked-go/repository"
)

const titleInUseErrorMessage = "title is already in use"

type IGalleryService interface {
	New(*model.NewGallery) (*model.Gallery, *model.Error)
}

type galleryService struct {
	repository.IGalleryRepository
}

func NewGalleryService(ur repository.IGalleryRepository) *galleryService {
	return &galleryService{
		ur,
	}
}

func (s *galleryService) New(create *model.NewGallery) (*model.Gallery, *model.Error) {
	create.Title = normalizeEmail(create.Title)

	gallery, _ := s.getByTitle(create.Title)

	if gallery != nil && gallery.UserId.String() == "TODO" {
		return nil, model.NewConflictApiError(titleInUseErrorMessage)
	}
	gallery = create.Gallery()

	err := s.Create(gallery)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *galleryService) getByTitle(title string) (*model.Gallery, *model.Error) {
	if title == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("title"))
	}
	gallery, err := s.Read("title", title)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

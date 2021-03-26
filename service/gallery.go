package service

import (
	"lens-locked-go/repository"
)

type IGalleryService interface {
}

type galleryService struct {
	repository.IGalleryRepository
}

func NewGalleryService(ur repository.IGalleryRepository) *galleryService {
	return &galleryService{
		ur,
	}
}

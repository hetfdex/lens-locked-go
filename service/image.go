package service

import (
	"bufio"
	"github.com/gofrs/uuid"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"os"
	"path/filepath"
)

type IImageService interface {
	Create(string, uuid.UUID) (*model.Image, *model.Error)
	GetById(uuid.UUID) (*model.Image, *model.Error)
	GetAllByGalleryId(uuid.UUID) ([]model.Image, *model.Error)
	Delete(*model.Image) *model.Error
}

type imageService struct {
	repository repository.IImageRepository
}

func NewImageService(ir repository.IImageRepository) *imageService {
	return &imageService{
		ir,
	}
}

func (s *imageService) Create(path string, galleryId uuid.UUID) (*model.Image, *model.Error) {
	image, err := makeImage(path, galleryId)

	if err != nil {
		return nil, err
	}
	err = s.repository.Create(image)

	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *imageService) GetById(id uuid.UUID) (*model.Image, *model.Error) {
	if id == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("id"))
	}
	gallery, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *imageService) GetAllByGalleryId(galleryId uuid.UUID) ([]model.Image, *model.Error) {
	if galleryId == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("galleryId"))
	}
	gallery, err := s.repository.ReadAll("gallery_id", galleryId)

	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *imageService) Delete(gallery *model.Image) *model.Error {
	return s.repository.Delete(gallery)
}

func makeImage(path string, galleryId uuid.UUID) (*model.Image, *model.Error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	defer file.Close()

	fileInfo, _ := file.Stat()

	fileSize := fileInfo.Size()

	bytes := make([]byte, fileSize)

	buffer := bufio.NewReader(file)

	_, err = buffer.Read(bytes)

	return &model.Image{
		Bytes:     bytes,
		Extension: filepath.Ext(path),
		GalleryId: galleryId,
	}, nil
}

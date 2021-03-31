package service

import (
	"bytes"
	"github.com/gofrs/uuid"
	"io"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"path/filepath"
	"strings"
)

const unsupportedFileErrorMessage = "unsupported file"

type IImageService interface {
	Create(io.ReadCloser, string, uuid.UUID) *model.Error
	GetById(uuid.UUID) (*model.Image, *model.Error)
	GetAllByGalleryId(uuid.UUID) ([]*model.Image, *model.Error)
}

type imageService struct {
	repository repository.IImageRepository
}

func NewImageService(ir repository.IImageRepository) *imageService {
	return &imageService{
		ir,
	}
}

func (s *imageService) Create(file io.ReadCloser, filename string, galleryId uuid.UUID) *model.Error {
	defer file.Close()

	filename = lower(filepath.Base(filename))

	extension := filepath.Ext(filename)

	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
		return model.NewBadRequestApiError(unsupportedFileErrorMessage)
	}
	buffer := bytes.NewBuffer(nil)

	_, err := io.Copy(buffer, file)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	filename = strings.TrimSuffix(filename, extension)

	image := &model.Image{
		Bytes:     buffer.Bytes(),
		Name:      filename,
		Extension: extension,
		GalleryId: galleryId,
	}

	er := s.repository.Create(image)

	if er != nil {
		return er
	}
	return nil
}

func (s *imageService) GetById(id uuid.UUID) (*model.Image, *model.Error) {
	if id == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("id"))
	}
	image, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *imageService) GetAllByGalleryId(galleryId uuid.UUID) ([]*model.Image, *model.Error) {
	if galleryId == uuid.Nil {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("galleryId"))
	}
	images, err := s.repository.ReadAll("gallery_id", galleryId)

	if err != nil {
		return nil, err
	}
	return images, nil
}

package service

import (
	"bytes"
	"github.com/gofrs/uuid"
	"io"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"mime/multipart"
	"path/filepath"
)

const unsupportedFileErrorMessage = "unsupported file"

type IImageService interface {
	Create(*multipart.FileHeader, uuid.UUID) *model.Error
	GetById(uuid.UUID) (*model.Image, *model.Error)
}

type imageService struct {
	repository repository.IImageRepository
}

func NewImageService(ir repository.IImageRepository) *imageService {
	return &imageService{
		ir,
	}
}

func (s *imageService) Create(fileHeader *multipart.FileHeader, galleryId uuid.UUID) *model.Error {
	extension := lower(filepath.Ext(fileHeader.Filename))

	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		file, err := fileHeader.Open()

		if err != nil {
			return model.NewInternalServerApiError(err.Error())
		}
		defer file.Close()

		buffer := bytes.NewBuffer(nil)

		_, err = io.Copy(buffer, file)

		image := &model.Image{
			Bytes:     buffer.Bytes(),
			Extension: extension,
			GalleryId: galleryId,
		}

		er := s.repository.Create(image)

		if er != nil {
			return er
		}
		return nil
	}
	return model.NewBadRequestApiError(unsupportedFileErrorMessage)
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

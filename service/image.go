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

	if extension != "jpg" && extension != "jpeg" && extension != "png" {
		return model.NewBadRequestApiError(unsupportedFileErrorMessage)
	}
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

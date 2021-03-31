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
	Create(*multipart.FileHeader, uuid.UUID) (*model.Image, *model.Error)
	GetById(uuid.UUID) (*model.Image, *model.Error)
	GetAllByGalleryId(uuid.UUID) ([]*model.Image, *model.Error)
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

func (s *imageService) Create(fileHeader *multipart.FileHeader, galleryId uuid.UUID) (*model.Image, *model.Error) {
	extension := lower(filepath.Ext(fileHeader.Filename))

	er := ValidExtension(extension)

	if er != nil {
		return nil, er
	}
	file, err := fileHeader.Open()

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	defer file.Close()

	buffer := bytes.NewBuffer(nil)

	_, err = io.Copy(buffer, file)

	image := &model.Image{
		Bytes:     buffer.Bytes(),
		Extension: extension,
		GalleryId: galleryId,
	}

	er = s.repository.Create(image)

	if er != nil {
		return nil, er
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

func (s *imageService) GetAllByGalleryId(galleryId uuid.UUID) ([]*model.Image, *model.Error) {
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

func ValidExtension(extension string) *model.Error {
	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		return nil
	}
	return model.NewBadRequestApiError(unsupportedFileErrorMessage)
}

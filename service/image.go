package service

import (
	"bufio"
	"encoding/base64"
	"github.com/gofrs/uuid"
	"io"
	"io/ioutil"
	"lens-locked-go/model"
	"lens-locked-go/repository"
	"lens-locked-go/util"
	"path/filepath"
	"strings"
)

const unsupportedFileErrorMessage = "unsupported file"

type IImageService interface {
	Create(io.ReadCloser, string, uuid.UUID) *model.Error
	GetById(uuid.UUID) (*model.Image, *model.Error)
	GetAllByGalleryId(uuid.UUID) ([]*model.Image, *model.Error)
	Delete(image *model.Image) *model.Error
}

type imageService struct {
	repository repository.IImageRepository
}

func newImageService(ir repository.IImageRepository) *imageService {
	return &imageService{
		ir,
	}
}

func (s *imageService) Create(file io.ReadCloser, filename string, galleryId uuid.UUID) *model.Error {
	defer file.Close()

	filename = lower(filename)

	extension := filepath.Ext(filename)

	filename = strings.TrimSuffix(filename, extension)

	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" {
		return model.NewBadRequestApiError(unsupportedFileErrorMessage)
	}
	reader := bufio.NewReader(file)

	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return model.NewInternalServerApiError(err.Error())
	}
	source := base64.StdEncoding.EncodeToString(content)

	image := &model.Image{
		Source:    source,
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
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("id"))
	}
	image, err := s.repository.Read("id", id)

	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *imageService) GetAllByGalleryId(galleryId uuid.UUID) ([]*model.Image, *model.Error) {
	if galleryId == uuid.Nil {
		return nil, model.NewInternalServerApiError(util.MustNotBeEmptyErrorMessage("galleryId"))
	}
	images, err := s.repository.ReadAll("gallery_id", galleryId)

	if err != nil {
		return nil, err
	}
	return images, nil
}

func (s *imageService) Delete(image *model.Image) *model.Error {
	return s.repository.Delete(image)
}

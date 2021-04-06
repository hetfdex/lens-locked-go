package service

import (
	"lens-locked-go/config"
	"lens-locked-go/repository"
)

type Services struct {
	User    IUserService
	Gallery IGalleryService
	Image   IImageService
	Dropbox IDropboxService
}

func NewServices(rp *repository.Repositories, cfg *config.Config) *Services {
	return &Services{
		User:    newUserService(rp.User, cfg.Crypto),
		Gallery: newGalleryService(rp.Gallery),
		Image:   newImageService(rp.Image),
		Dropbox: newDropboxService(rp.Dropbox),
	}
}

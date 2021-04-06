package repository

import "gorm.io/gorm"

type Repositories struct {
	User    IUserRepository
	Gallery IGalleryRepository
	Image   IImageRepository
	Dropbox IDropboxRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:    newUserRepository(db),
		Gallery: newGalleryRepository(db),
		Image:   newImageRepository(db),
		Dropbox: newDropboxRepository(db),
	}
}

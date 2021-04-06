package wrapper

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxf "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"lens-locked-go/model"
)

type File struct {
	Name string
	Path string
}

type Folder struct {
	Name string
	Path string
}

func ListAll(accessToken string, path string) ([]*File, []*Folder, *model.Error) {
	dbxCfg := dropbox.Config{
		Token: accessToken,
	}
	client := dbxf.New(dbxCfg)

	arg := &dbxf.ListFolderArg{
		Path: path,
	}

	results, err := client.ListFolder(arg)

	if err != nil {
		return nil, nil, model.NewFailedDependencyApiError(err.Error())
	}

	var files []*File
	var folders []*Folder

	for _, entry := range results.Entries {
		switch meta := entry.(type) {
		case *dbxf.FileMetadata:
			file := &File{
				Name: meta.Name,
				Path: meta.PathLower,
			}

			files = append(files, file)
		case *dbxf.FolderMetadata:
			folder := &Folder{
				Name: meta.Name,
				Path: meta.PathLower,
			}

			folders = append(folders, folder)
		}
	}
	return files, folders, nil
}

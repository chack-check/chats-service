package generic_factories

import (
	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
)

func DbFileToSchema(dbFile models.SavedFile) model.SavedFile {
	return model.SavedFile{
		OriginalURL:       dbFile.OriginalUrl,
		OriginalFilename:  dbFile.OriginalFilename,
		ConvertedURL:      &dbFile.ConvertedUrl,
		ConvertedFilename: &dbFile.ConvertedFilename,
	}
}

func SchemaSavedFileToDbFile(schemaFile model.SavedFile) models.SavedFile {
	if schemaFile.ConvertedURL != nil && schemaFile.ConvertedFilename != nil {
		return models.SavedFile{
			OriginalUrl:       schemaFile.OriginalURL,
			OriginalFilename:  schemaFile.OriginalFilename,
			ConvertedUrl:      *schemaFile.ConvertedURL,
			ConvertedFilename: *schemaFile.ConvertedFilename,
		}
	} else {
		return models.SavedFile{
			OriginalUrl:      schemaFile.OriginalURL,
			OriginalFilename: schemaFile.OriginalFilename,
		}
	}
}

func UploadingFileToDbFile(uploadingFile model.UploadingFile) models.SavedFile {
	if uploadingFile.Converted != nil {
		return models.SavedFile{
			OriginalUrl:       uploadingFile.Original.URL,
			OriginalFilename:  uploadingFile.Original.Filename,
			ConvertedUrl:      uploadingFile.Converted.URL,
			ConvertedFilename: uploadingFile.Converted.Filename,
		}
	} else {
		return models.SavedFile{
			OriginalUrl:      uploadingFile.Original.URL,
			OriginalFilename: uploadingFile.Original.Filename,
		}
	}
}

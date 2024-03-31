package factories

import (
	"github.com/chack-check/chats-service/api/v1/dtos"
	"github.com/chack-check/chats-service/api/v1/graph/model"
	"github.com/chack-check/chats-service/api/v1/models"
)

func SavedFileToFileDto(file models.SavedFile) dtos.FileDto {
	var converted_url *string
	var converted_filename *string
	if file.ConvertedUrl == "" {
		converted_url = nil
	} else {
		converted_url = &file.ConvertedUrl
	}

	if file.ConvertedFilename == "" {
		converted_filename = nil
	} else {
		converted_filename = &file.ConvertedFilename
	}

	return dtos.FileDto{
		OriginalUrl:       file.OriginalUrl,
		OriginalFilename:  file.OriginalFilename,
		ConvertedUrl:      converted_url,
		ConvertedFilename: converted_filename,
	}
}

func FileDtoToSavedFile(file dtos.FileDto) models.SavedFile {
	var converted_url string
	var converted_filename string
	if file.ConvertedUrl != nil {
		converted_url = *file.ConvertedUrl
	} else {
		converted_url = ""
	}

	if file.ConvertedFilename != nil {
		converted_filename = *file.ConvertedFilename
	} else {
		converted_filename = ""
	}

	return models.SavedFile{
		OriginalUrl:       file.OriginalUrl,
		OriginalFilename:  file.OriginalFilename,
		ConvertedUrl:      converted_url,
		ConvertedFilename: converted_filename,
	}
}

func FileDtoToSchema(dbFile dtos.FileDto) model.SavedFile {
	return model.SavedFile{
		OriginalURL:       dbFile.OriginalUrl,
		OriginalFilename:  dbFile.OriginalFilename,
		ConvertedURL:      dbFile.ConvertedUrl,
		ConvertedFilename: dbFile.ConvertedFilename,
	}
}

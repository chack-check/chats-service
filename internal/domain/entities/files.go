package entities

type SavedFile struct {
	originalURL       string
	originalFilename  string
	convertedURL      *string
	convertedFilename *string
}

func NewSavedFile(originalurl, originalFilename string, convertedURL, convertedFilename *string) SavedFile {
	return SavedFile{
		originalURL:       originalurl,
		originalFilename:  originalFilename,
		convertedURL:      convertedURL,
		convertedFilename: convertedFilename,
	}
}

func (model *SavedFile) GetOriginalURL() string {
	return model.originalURL
}

func (model *SavedFile) GetOriginalFilename() string {
	return model.originalFilename
}

func (model *SavedFile) GetConvertedURL() *string {
	return model.convertedURL
}

func (model *SavedFile) GetConvertedFilename() *string {
	return model.convertedFilename
}

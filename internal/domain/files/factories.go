package files

func UploadingFileToSavedFile(file UploadingFile) SavedFile {
	var convertedUrl *string
	var convertedFilename *string
	if file.GetConverted() != nil {
		url := file.GetConverted().GetUrl()
		filename := file.GetConverted().GetFilename()
		convertedUrl = &url
		convertedFilename = &filename
	} else {
		convertedUrl = nil
		convertedFilename = nil
	}

	return NewSavedFile(
		file.GetOriginal().GetUrl(),
		file.GetOriginal().GetFilename(),
		convertedUrl,
		convertedFilename,
	)
}

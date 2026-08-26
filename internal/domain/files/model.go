package files

type SavedFile struct {
	originalUrl       string
	originalFilename  string
	convertedUrl      *string
	convertedFilename *string
}

func (model *SavedFile) GetOriginalUrl() string {
	return model.originalUrl
}

func (model *SavedFile) GetOriginalFilename() string {
	return model.originalFilename
}

func (model *SavedFile) GetConvertedUrl() *string {
	return model.convertedUrl
}

func (model *SavedFile) GetConvertedFilename() *string {
	return model.convertedFilename
}

type SystemFiletype string

func (e SystemFiletype) IsValid() bool {
	switch e {
	case AvatarFiletype, FileInChatFiletype, VoiceFiletype, CircleFiletype:
		return true
	}

	return false
}

func (e SystemFiletype) String() string {
	return string(e)
}

const (
	AvatarFiletype     SystemFiletype = "avatar"
	FileInChatFiletype SystemFiletype = "file_in_chat"
	VoiceFiletype      SystemFiletype = "voice"
	CircleFiletype     SystemFiletype = "circle"
)

type UploadingFileMeta struct {
	url            string
	filename       string
	signature      string
	systemFiletype SystemFiletype
}

func (model UploadingFileMeta) GetUrl() string {
	return model.url
}

func (model UploadingFileMeta) GetFilename() string {
	return model.filename
}

func (model UploadingFileMeta) GetSignature() string {
	return model.signature
}

func (model UploadingFileMeta) GetSystemFiletype() SystemFiletype {
	return model.systemFiletype
}

type UploadingFile struct {
	original  UploadingFileMeta
	converted *UploadingFileMeta
}

func (model *UploadingFile) GetOriginal() UploadingFileMeta {
	return model.original
}

func (model *UploadingFile) GetConverted() *UploadingFileMeta {
	return model.converted
}

func NewSavedFile(originalurl, originalFilename string, convertedUrl, convertedFilename *string) SavedFile {
	return SavedFile{
		originalUrl:       originalurl,
		originalFilename:  originalFilename,
		convertedUrl:      convertedUrl,
		convertedFilename: convertedFilename,
	}
}

func NewUploadingFileMeta(url, filename, signature string, systemFiletype SystemFiletype) UploadingFileMeta {
	return UploadingFileMeta{
		url:            url,
		filename:       filename,
		signature:      signature,
		systemFiletype: systemFiletype,
	}
}

func NewUploadingFile(original UploadingFileMeta, converted *UploadingFileMeta) UploadingFile {
	return UploadingFile{
		original:  original,
		converted: converted,
	}
}

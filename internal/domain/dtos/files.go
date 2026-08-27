package dtos

import "chats-service/internal/domain/constants"

type UploadingFileMeta struct {
	url            string
	filename       string
	signature      string
	systemFiletype constants.SystemFiletype
}

func NewUploadingFileMeta(url, filename, signature string, systemFiletype constants.SystemFiletype) UploadingFileMeta {
	return UploadingFileMeta{
		url:            url,
		filename:       filename,
		signature:      signature,
		systemFiletype: systemFiletype,
	}
}

func (model UploadingFileMeta) GetURL() string {
	return model.url
}

func (model UploadingFileMeta) GetFilename() string {
	return model.filename
}

func (model UploadingFileMeta) GetSignature() string {
	return model.signature
}

func (model UploadingFileMeta) GetSystemFiletype() constants.SystemFiletype {
	return model.systemFiletype
}

type UploadingFile struct {
	original  UploadingFileMeta
	converted *UploadingFileMeta
}

func NewUploadingFile(original UploadingFileMeta, converted *UploadingFileMeta) UploadingFile {
	return UploadingFile{
		original:  original,
		converted: converted,
	}
}

func (model *UploadingFile) GetOriginal() UploadingFileMeta {
	return model.original
}

func (model *UploadingFile) GetConverted() *UploadingFileMeta {
	return model.converted
}

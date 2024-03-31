package dtos

type FileDto struct {
	OriginalUrl       string  `json:"original_url"`
	OriginalFilename  string  `json:"original_filename"`
	ConvertedUrl      *string `json:"converted_url"`
	ConvertedFilename *string `json:"converted_filename"`
}

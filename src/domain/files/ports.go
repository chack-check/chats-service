package files

type FilesPort interface {
	GetSignatureForFile(filename string, systemFiletype SystemFiletype) string
}

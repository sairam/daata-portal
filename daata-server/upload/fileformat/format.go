package fileformat

//FileFormat is a should be used for understand file extension. we should be able to add more like bz2, gz etc.,
type FileFormat int

// Support file formats
const (
	FileText FileFormat = iota + 1
	FileJSON
	FileHTML
	FileZip
	FileGz
	FileTar
)

package fileformat

import (
	"regexp"
	"strings"
)

type CompressionFormat int

const (
	CompressionNone CompressionFormat = iota + 1
	CompressionBz2
	CompressionGz
)

func FindCompressionFormat(str string) (string, CompressionFormat) {
	var c CompressionFormat
	// bz2, gz
	var extRegexp = regexp.MustCompile(`(gz|bz2)$`)
	var match = extRegexp.FindStringSubmatch(str)
	if len(match) > 0 {
		switch match[0] {
		case "bz2":
			c = CompressionBz2
		case "gz":
			c = CompressionGz
		default:
			c = CompressionNone
		}
	}
	if len(match) == 0 {
		return str, CompressionNone
	}
	return strings.Replace(str, "."+match[0], "", 1), c
}

type ArchiveFormat int

const (
	ArchiveNone ArchiveFormat = iota + 1
	ArchiveZip
	ArchiveTar
	// ArchiveZ
)

func FindArchiveFormat(str string) (string, ArchiveFormat) {
	var c ArchiveFormat
	var extRegexp = regexp.MustCompile(`(zip|tar)$`)
	var match = extRegexp.FindStringSubmatch(str)
	if len(match) > 0 {
		switch match[0] {
		case "tar":
			c = ArchiveTar
		case "zip":
			c = ArchiveZip
		default:
			c = ArchiveNone
		}
	}
	if len(match) == 0 {
		return str, ArchiveNone
	}
	return strings.Replace(str, "."+match[0], "", 1), c
}

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

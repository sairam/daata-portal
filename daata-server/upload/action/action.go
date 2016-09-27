package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	ff "../fileformat"
	"./gz"
	"./tar"
	"./zip"
)

// Settings are required to perform actions in the order we require
type Settings struct {
	CompressionType ff.CompressionFormat // tar.bz2 - first its uncompressed
	ArchiveType     ff.ArchiveFormat     // second, its untarred or unzipped
	AppendMode      bool                 // append mode or override file
}

// Perform goes to a directory, performs the action based on the format
// and returns to the current directory
func Perform(settings *Settings, dir, file string) string {

	currentDirectory, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(currentDirectory)
	location := "."

	file = decompress(settings, file, location)
	dirPath, requiresCleanup := unarchive(settings, file, location)
	if requiresCleanup {
		cleanup(dirPath, file)
	}

	return ""
}

func unarchive(settings *Settings, file, location string) (string, bool) {
	var err error
	var out []byte
	var cleanup = true

	switch settings.ArchiveType {
	case ff.ArchiveZip:
		out, err = zip.Extract(file, location)
	case ff.ArchiveTar:
		out, err = tar.Extract(file, location)
	default:
		out, err = []byte{}, nil
		cleanup = false
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output is \n%s\n", out)

	// TODO - remove archived file
	// return string(out)
	return location, cleanup

}

func decompress(settings *Settings, file, location string) string {
	var newFile, oldFile string
	var err error

	switch settings.CompressionType {
	case ff.CompressionGz:
		newFile, oldFile, err = gz.Decompress(file, location)
	case ff.CompressionBz2:
		// TODO , fix me
		newFile, oldFile, err = gz.Decompress(file, location)
	default:
		newFile, oldFile, err = file, "", nil
	}
	if err != nil {
		fmt.Println(err)
	}
	if oldFile != "" {
		os.Remove(oldFile)
	}
	// TODO - remove oldFile if present
	return newFile
}

// Cleanup zip files
func cleanup(dir, file string) {
	os.Remove(file)
	files, _ := ioutil.ReadDir("./")
	if len(files) == 1 {
		file := files[0]
		if file.IsDir() {
			moveFilesToParent(file.Name())
		}
	}
}

// move files from single sub directory to current directory
func moveFilesToParent(dir string) error {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		filename := strings.Join([]string{dir, file.Name()}, string(os.PathSeparator))
		os.Rename(filename, file.Name())
	}
	os.Remove(dir)
	return nil
}

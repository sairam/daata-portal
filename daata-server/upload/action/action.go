package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	f "../fileformat"
	_ "./gz"  // f
	_ "./tar" // f
	"./zip"
)

// Perform goes to a directory, performs the action based on the format
// and returns to the current directory
func Perform(format f.FileFormat, dir, file string) string {
	currentDirectory, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(currentDirectory)
	if format == f.FileZip {
		location := "."
		out, err := zip.Extract(file, location)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Output is \n%s\n", out)
		cleanup(dir, file)
		return string(out)
	}
	return ""
}
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

func moveFilesToParent(dir string) error {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		filename := strings.Join([]string{dir, file.Name()}, string(os.PathSeparator))
		os.Rename(filename, file.Name())
	}
	os.Remove(dir)
	return nil
}

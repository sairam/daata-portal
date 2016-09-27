package action

import (
	"fmt"
	"log"
	"os"

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
		return string(out)
	}
	return ""
}

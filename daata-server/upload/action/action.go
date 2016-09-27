package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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
		cmd := []string{"/bin/mv", files[0].Name() + "/*", "./"}
		fmt.Printf("No of files are %s\n", files[0].Name())
		fmt.Println(cmd)
		exec.Command(cmd[0], cmd[1:]...).Output()
	}
	fmt.Printf("No of files are %d\n", len(files))
	for file := range files {
		fmt.Println(file)
	}
}

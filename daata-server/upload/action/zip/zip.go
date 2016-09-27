package zip

import (
	"fmt"
	"os"
	"os/exec"
)

// Interface for action
// type Interface interface {
// 	Extract() ([]byte, error)
// 	List() ([]byte, error)
// }

// ExtractHere extracts the file in the current directory.
// Also provide output on the details of the files
func ExtractHere(file string) ([]byte, error) {
	return Extract(file, ".")
}

// TODO -take care of dot files which are generated

// Extract file to a particular location
// Also provide output on the details of the files
func Extract(file, location string) ([]byte, error) {
	fmt.Println(os.Getwd())
	cmd := []string{"/usr/bin/unzip", file, "-d", location}
	fmt.Println(cmd)
	return exec.Command(cmd[0], cmd[1:]...).Output()
}

// List the files in the zip
func List(file string) ([]byte, error) {
	cmd := []string{"/usr/bin/unzip", "-l", file}
	return exec.Command(cmd[0], cmd[1:]...).Output()
}

package bz2

import (
	"fmt"
	"os"
	"os/exec"
)

// ExtractHere extracts the file in the current directory.
// Also provide output on the details of the files
func ExtractHere(file string) ([]byte, error) {
	return Extract(file, ".")
}

// TODO -take care of dot files which are generated

// Extract file to a particular location
// Also provide output on the details of the files
func Extract(file, location string) ([]byte, error) {
	fmt.Println("Location is ignored")
	fmt.Println(os.Getwd())
	cmd := []string{"/usr/bin/bunzip2", file}
	fmt.Println(cmd)

	// Display everything we got if error.
	output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		fmt.Println("Error when running command.  Output:")
		fmt.Println(string(output))
		fmt.Printf("Got command status: %s\n", err.Error())
	}
	return output, err
	// return exec.Command(cmd[0], cmd[1:]...).Output()
}

// List the files in the zip
func List(file string) ([]byte, error) {
	cmd := []string{"/usr/bin/bunzip2", "-l", file}
	return exec.Command(cmd[0], cmd[1:]...).Output()
}

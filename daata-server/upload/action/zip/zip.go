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

// // Display just the stderr if an error occurs
// cmd := exec.Command(...)
// stderr := &bytes.Buffer{}    // make sure to import bytes
// cmd.Stderr = stderr
// err := cmd.Run()
// if err != nil {
//     fmt.Println("Error when running command.  Error log:")
//     fmt.Println(stderr.String())
//     fmt.Printf("Got command status: %s\n", err.Error())
//     return
// }

// List the files in the zip
func List(file string) ([]byte, error) {
	cmd := []string{"/usr/bin/unzip", "-l", file}
	return exec.Command(cmd[0], cmd[1:]...).Output()
}

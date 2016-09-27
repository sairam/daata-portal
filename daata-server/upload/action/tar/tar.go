package tar

import "os/exec"

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

// Extract file to a particular location
// FIXME
// Also provide output on the details of the files
func Extract(file, location string) ([]byte, error) {
	cmd := []string{"/usr/bin/tar", "-cf", file, "-d", location}
	return exec.Command(cmd[0], cmd[1:]...).Output()
}

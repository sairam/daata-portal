package gz

import "os/exec"

// Interface for action
// type Interface interface {
// 	Extract() ([]byte, error)
// 	List() ([]byte, error)
// }

// ExtractHere extracts the file in the current directory.
// Also provide output on the details of the files
func ExtractHere(file string) (string, string, error) {
	return Decompress(file, ".")
}

// Extract file to a particular location
// FIXME
// Also provide output on the details of the files
func Decompress(file, location string) (string, string, error) {
	cmd := []string{"/usr/bin/gunzip", file, "-d", location}
	exec.Command(cmd[0], cmd[1:]...).Output()
	return "", "", nil
}

package upload

import (
	"os"
	"strings"

	"../config"
)

type uploadLocation struct {
	directory string
	subdir    string
	aliases   []string
}

// clean ensure directory starts with a string instead of a /
// subdirectory should not have any path separators???
// aliases do not have path separators???
func (u *uploadLocation) clean() error {
	return nil
}

// Generate the directory structure
func (u *uploadLocation) generateDirectory() error {
	err := os.MkdirAll(u.directory, config.DirectoryPermissions)
	if err != nil {
		return err
	}
	err = os.Mkdir(u.path(), config.DirectoryPermissions)
	return err
}

// Give the path of subdirectory w/ directory. Essentially the complete OS path relative
func (u *uploadLocation) path() string {
	return strings.Join([]string{u.directory, u.subdir}, string(os.PathSeparator))
}

// Give the path of any location w/ directory. Used for making aliases to directories
func (u *uploadLocation) abspath(location string) string {
	return strings.Join([]string{u.directory, location}, string(os.PathSeparator))
}

// aliases are made from the newly created subdirectory to newer locations for easier linking
func (u *uploadLocation) makeAliases() []error {
	var err []error
	for _, alias := range u.aliases {
		lerr := os.Symlink(u.path(), u.abspath(alias))
		if lerr != nil {
			err = append(err, lerr)
		}
	}
	return err
}

package upload

import (
	"os"
	"strings"

	conf "../config"
)

type uploadLocation struct {
	directory string
	subdir    string
	aliases   []string
	filename  string
	extension string // extension is to be appended to filename
}

// clean ensure directory starts with a string instead of a /
// subdirectory should not have any path separators???
// aliases do not have path separators???
func (u *uploadLocation) clean() error {
	return nil
}

// Generate the directory structure
func (u *uploadLocation) generateDirectory() error {
	// fmt.Println("created directory at ", u.dirpath())
	err := os.MkdirAll(u.dirpath(), os.FileMode(conf.C().Permissions.Directory))
	return err
}

// Give the path of subdirectory w/ directory. Essentially the complete OS path relative
func (u *uploadLocation) dirpath() string {
	if u.subdir == "" {
		return u.directory
	}
	return strings.Join([]string{u.directory, u.subdir}, string(os.PathSeparator))
}

func (u *uploadLocation) path() string {
	return strings.Join([]string{u.dirpath(), u.filepath()}, string(os.PathSeparator))
}

func (u *uploadLocation) filepath() string {
	if u.filename == "" {
		return ""
	}
	return u.filename + "." + u.extension
}

// Give the path of any location w/ directory. Used for making aliases to directories
func (u *uploadLocation) abspath(location string) string {
	return strings.Join([]string{u.directory, location}, string(os.PathSeparator))
}

// aliases are made from the newly created subdirectory to newer locations for easier linking
// Aliases cannot be made to actual files
func (u *uploadLocation) makeAliases() []error {
	var err []error
	// Aliases are not possible when subdirs are not present
	if u.subdir == "" {
		return err
	}
	// fmt.Printf("Aliases are %v\n", u.aliases)

	for _, alias := range u.aliases {
		os.Remove(alias)
		lerr := os.Symlink(u.subdir, alias)
		if lerr != nil {
			// fmt.Println(lerr)
			err = append(err, lerr)
		}
	}
	return err
}

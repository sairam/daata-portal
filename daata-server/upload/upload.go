package upload

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type uploadType int

const (
	undetermined uploadType = iota
	static                  // directory/file
	versioned               // directory/file
	dataPoint               // inside a directory. each data point is a file
	table                   // inside a directory, its a file
	parallel                // per file
)

/*
  upload.go determines based on the parameters or existing directory structure
  which functions to call.

  1. Static files which are versioned and/or aliased identified by Type: 'VD'
  2. Static files which are just uploaded
  3. Data sent in form of key/value for graphs
  4. Static files edited in UI (via markdown etc., to be updated in place)
*/

func unableToDetermine(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Unable to determine the upload type. Internal Server Error!")
}

// TODO use a struct to load the type into memory
type uploadSettings struct {
	uploadType uploadType
	http.ResponseWriter
	*http.Request
}

func registerUpload() {
}

func determineUploadType(w http.ResponseWriter, r *http.Request) {

	function := unableToDetermine
	// settings := &uploadSettings{undetermined, w, r}
	switch strings.Join(r.Header["X_Upload_Type"], "") {
	case "static", "Static", "onetime", "OneTime":
		function = UploadStatic // static files like zip or html without versioning (below code to SaveFile)
	case "versioned_files", "VersionedFiles":
		function = UploadVersioned
	case "data_point", "data_points", "dataPoint", "dataPoints":
		function = UploadDataPoints // data points like key/value one or multiple
	case "table", "Table":
		function = UploadTable // json or CSV formats
	case "parallel", "Parallel":
		function = UploadParallel // multiple files to be put into the same location appended

	default:
		function = unableToDetermine

	}
	_ = function
	function(w, r)
}

/*
Later
Access restriction
If unzip file does not contain index.html, generate one with a tree.
*/

//FileFormat is a should be used for understand file extension
type FileFormat int

const (
	text FileFormat = iota + 1
	json
	html
	zip
)

func performAction(format FileFormat, dir, file string) string {
	currentDirectory, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(currentDirectory)
	if format == zip {
		cmd := []string{"/usr/bin/unzip", file}
		out, err := exec.Command(cmd[0], cmd[1]).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Output is \n%s\n", out)
		return string(out)
	}
	return ""
}

func getAction(ext string) FileFormat {
	if ext == "zip" {
		return zip
	}
	return text
}

func init() {
	http.HandleFunc("/u/", determineUploadType)
}

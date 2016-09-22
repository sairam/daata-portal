package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
  upload.go determines based on the parameters or existing directory structure
  which functions to call.

  1. Static files which are versioned and/or aliased identified by Type: 'VD'
  2. Static files which are just uploaded
  3. Data sent in form of key/value for graphs
  4. Static files edited in UI (via markdown etc., to be updated in place)
*/

func unableToDetermine(_ http.ResponseWriter, _ *http.Request) error {
	return errors.New("Unable to determine the upload type. Internal Server Error!")
}

func determineUploadType() {

	function := unableToDetermine

	switch r.Header["X_Upload_Type"] {
	case "static", "Static", "onetime", "OneTime":
		function = "upload_static" // static files like zip or html without versioning (below code to SaveFile)
	case "versioned_files", "VersionedFiles":
		function = upload_versioned
	case "data_point", "data_points", "dataPoint", "dataPoints":
		function = "upload_data" // data points like key/value one or multiple
	case "table", "Table":
		function = "upload_table" // json or CSV formats
	case "parallel", "Parallel":
		function = "upload_parallel" // multiple files to be put into the same location appended

	default:
		function = unableToDetermine

	}
	return function()
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	// 0. generate random id
	dirName := randomString(randomStringLength)
	url := serverURL + "/d/" + dirName
	// 1. read contents
	data, _ := ioutil.ReadAll(r.Body)
	debug(w, r)

	// 2. save file
	extension := strings.Split(r.Header["Content-Type"][0], "/")[1]
	directory, fileName := saveToFile(dirName, extension, data)

	// 3. determine file type
	action := getAction(extension)

	// 4. perform action of unzip or nothing
	// 5. TODO - add symlinks as per provided option
	output := performAction(action, directory, fileName)
	fmt.Fprintf(w, "\n"+output+"\n")

	// 5. send back url based on random id
	fmt.Fprintf(w, "\n"+url+"\n")
}

/*
Later
Access restriction
If unzip file does not contain index.html, generate one with a tree.
*/

func init() {
	http.HandleFunc("/u/", saveFile)
}

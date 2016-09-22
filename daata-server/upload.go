package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
  upload.go determines based on the parameters or existing directory structure
  which functions to call.
*/

//
// func showFile(w http.ResponseWriter, r *http.Request) {
// 	fileName := r.URL.String()
// 	data := strings.TrimLeft(fileName, "/d/")
// 	openIndex(fileName)
// 	fmt.Fprintf(w, "%s\n", data)
// }
//
// func openIndex(file string) {
//
// }

// "http://localhost" + port
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

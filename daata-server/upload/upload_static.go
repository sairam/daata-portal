package upload

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"../config"
	"../utils"
)

/*
  uploading static files takes only static files hosted under a directory.
  If a directory is not specified, one will be allocated to it.

  Future callers should call that directory
*/

// UploadStatic files
func UploadStatic(w http.ResponseWriter, r *http.Request) error {
	// 0. generate random id
	dirName := utils.RandomString(config.RandomStringLength)
	url := config.ServerURL + "/d/" + dirName
	// 1. read contents
	data, _ := ioutil.ReadAll(r.Body)
	utils.DebugHTTP(w, r)

	// 2. save file
	extension := strings.Split(r.Header["Content-Type"][0], "/")[1]
	directory, fileName := utils.SaveToFile(dirName, extension, data)

	// 3. determine file type
	action := getAction(extension)

	// 4. perform action of unzip or nothing
	// 5. TODO - add symlinks as per provided option
	output := performAction(action, directory, fileName)
	fmt.Fprintf(w, "\n"+output+"\n")

	// 5. send back url based on random id
	fmt.Fprintf(w, "\n"+url+"\n")

	return nil
}

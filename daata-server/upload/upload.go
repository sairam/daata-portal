package upload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"../config"
	"../utils"
)

/*
  upload.go determines based on the parameters or existing directory structure
  which functions to call.

  1. Static files which are versioned and/or aliased identified by Type: 'VD'
  2. Static files which are just uploaded
  3. Data sent in form of key/value for graphing
  4. Static files edited in UI (via markdown etc., to be updated in place)

  Static:
  curl -i -X POST -H "Content-Type: application/zip" --data-binary "@data.zip" https://my.daata.xyz/u/docs/spokes-platform

  Versioned:
  curl -X POST -H 'X-Version: "2.1.3"' -H 'X-Alias: release-20160707, master, stable' -H 'Content-Type: application/zip' --file-binary="@filename.zip" https://my.daata.xyz/docs/spokes-platform

  fetching the data, attributes based on headers, file extension type, parsing, etc.,

*/

type uploadType int

const (
	undetermined uploadType = iota
	static                  // directory/file
	versioned               // directory/file
	dataPoint               // inside a directory. each data point is a file
	table                   // inside a directory, its a file
	parallel                // per file
)

//FileFormat is a should be used for understand file extension. we should be able to add more like bz2, gz etc.,
type fileFormat int

const (
	text fileFormat = iota + 1
	json
	html
	zip
)

type upload struct{}

type uploadLocation struct {
	directory string
	subdir    string
	aliases   []string
}

func (u *uploadLocation) generateDirectory() error {
	err := os.MkdirAll(u.directory, config.DirectoryPermissions)
	if err != nil {
		return err
	}
	err = os.Mkdir(u.path(), config.DirectoryPermissions)
	return err
}

func (u *uploadLocation) path() string {
	return strings.Join([]string{u.directory, u.subdir}, string(os.PathSeparator))
}

func (u *uploadLocation) makeAliases() []error {
	var err []error
	for _, alias := range u.aliases {
		lerr := os.Symlink(u.subdir, alias)
		if lerr != nil {
			err = append(err, lerr)
		}
	}
	return err
}

// 1. determine actions to do
// 2. call the relevant functions
/*
  1. Get the path to upload the directory to
  2. Determine directory path based on other headers
  3. Get the extension of the file uploaded
  4. Based on the extension/request format, determine how to take action on the contents
*/
func unableToDetermine(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Unable to determine the upload type. Internal Server Error!")
}

func (u *upload) delegate(w http.ResponseWriter, r *http.Request) {

	function := unableToDetermine
	// settings := &uploadSettings{undetermined, w, r}
	switch strings.Join(r.Header["X_Upload_Type"], "") {
	case "static", "Static", "onetime", "OneTime":
		function = Static // static files like zip or html without versioning (below code to SaveFile)
	case "versioned_files", "VersionedFiles":
		function = Versioned
	case "data_point", "data_points", "dataPoint", "dataPoints":
		function = DataPoints // data points like key/value one or multiple
	case "table", "Table":
		function = Table // json or CSV formats
	case "parallel", "Parallel":
		function = Parallel // multiple files to be put into the same location appended

	default:
		function = unableToDetermine

	}
	/*
	  1.1. Get path
	  1.2. If path is not there, create a directory under level 2 like if uploaded to /u/docs, create /u/docs/{123456}
	  1.3. If path is not there, create a directory under level 1 like if uploaded to /u, create /u/{123456}
	*/
	// Determine location to upload
	dirName, ok := getFromPath(r.URL.Path)
	if !ok {
		dirName = dirName + "/" + utils.RandomString(config.RandomStringLength)
	}
	dirName = convertPathToDirectory(dirName)

	subDirectory := getSubDirectory(r.Header)
	softLinks := getSoftLinks(r.Header)

	uploadLoc := &uploadLocation{dirName, subDirectory, softLinks}
	uploadLoc.generateDirectory()
	targetURL := config.ServerURL + "/d/" + convertDirectoryToPath(uploadLoc.path())
	w.Header().Add("X-Generated-URL", targetURL)

	// . determine file type uploaded
	extension := extBasedOnContentType(r.Header["Content-Type"])
	action := getAction(extension)

	// . Determine subdirectory to store based on other information
	switch action {
	case zip:
		// save single file
		// then unzip
		// remove/save source file
	case html, text, json:
		// save single file
		// TODO - add for form submit type if json or encoding etc.,
	}

	// read the data from the body based on type
	data, _ := ioutil.ReadAll(r.Body)
	utils.DebugHTTP(w, r)

	directory, fileName := utils.SaveToFile(dirName, extension, data)
	output := performAction(action, directory, fileName)
	fmt.Fprintf(w, "\n"+output+"\n")

	_ = output

	function(w, r)
}

func convertPathToDirectory(path string) string {
	str := strings.Split(path, "/")
	return strings.Join(str, string(os.PathSeparator))
}

func convertDirectoryToPath(path string) string {
	str := strings.Split(path, string(os.PathSeparator))
	return strings.Join(str, "/")
}

func getSubDirectory(header http.Header) string {
	return header["X-Version"][0]
}

func getSoftLinks(header http.Header) []string {
	return strings.Split(header["X-Alias"][0], ",")
}

func getFromPath(path string) (string, bool) {
	// strip out "/u/", then split by "/" to see the size
	newPath := strings.TrimLeft(path, UploadPrefix+"/")

	// remove blank strings
	data := strings.Split(newPath, "/")

	data = cleanStrings(data, "")
	return strings.Join(data, "/"), (len(data) >= 2)
}

func cleanStrings(data []string, selector string) []string {
	var r []string
	for _, str := range data {
		if str != selector {
			r = append(r, str)
		}
	}
	return r
}

func extBasedOnContentType(contentType []string) string {
	return strings.Split(contentType[0], "/")[1]
}

func getAction(ext string) fileFormat {
	switch ext {
	case "zip":
		return zip
	case "json":
		return json
	case "html":
		return html
	case "plain":
		return text
	default:
		return text
	}
}

func performAction(format fileFormat, dir, file string) string {
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

//UploadPrefix is required to get the upload prefix
const UploadPrefix = "/u"

func prefix() string {
	return UploadPrefix + "/"
}

func init() {
	u := &upload{}
	http.HandleFunc(prefix(), u.delegate)
}

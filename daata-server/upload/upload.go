package upload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"../config"
	"../display"
	"../utils"
	"./action"
	ff "./fileformat/"
)

// "/" is being used at many locations for many operations. A Path separator makes sense since we have a file separator
// We need this if we want to run the server on cross platform later on.
const httpPathSeparator = "/"

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
	// Determine location to upload
	uploadLoc := generateUploadLocation(r)

	// Later - Need to fork off here based on form vs upload type
	// If form data is submitted, we need to do action per file

	fileName := getFilename(r, getExt(r))
	theaction := getAction(fileName)

	// read the data from the body based on type
	data, _ := ioutil.ReadAll(r.Body)
	utils.DebugHTTP(w, r)

	// get main directory to save directory
	uploadLoc.generateDirectory()

	// save file in directory location
	directory, fileName := utils.SaveToFile(uploadLoc.path(), fileName, data)
	output := action.Perform(theaction, directory, fileName)
	fmt.Fprintf(w, "\n"+output+"\n")

	uploadLoc.makeAliases()

	targetURL := config.ServerURL + display.Prefix() + convertDirectoryToPath(uploadLoc.path())
	w.Header().Add("X-Generated-URL", targetURL)

	function(w, r)
}

func getFilename(r *http.Request, ext string) string {
	return "samplefilename"
}
func getExt(r *http.Request) string {
	// d := strings.Split(r.URL.Path, httpPathSeparator)
	// lastPath := d[len(d)-1]
	// TODO - regexp match
	// "(.tar.gz|.gz|.zip|.bz2|)$"
	// if strings.Contains(lastPath, )
	return extBasedOnContentType(r.Header["Content-Type"])
}

func extBasedOnContentType(contentType []string) string {
	return strings.Split(contentType[0], httpPathSeparator)[1]
}

func generateUploadLocation(r *http.Request) *uploadLocation {
	/*
	  1.1. Get path
	  1.2. If path is not there, create a directory under level 2 like if uploaded to /u/docs, create /u/docs/{123456}/
	  1.3. If path is not there, create a directory under level 1 like if uploaded to /u, create /u/{123456}/
	*/
	dirName, ok := getFromPath(r.URL.Path)
	if !ok {
		dirName = dirName + httpPathSeparator + utils.RandomString(config.RandomStringLength)
	}
	dirName = convertPathToDirectory(dirName)
	subDirectory := getSubDirectory(r.Header)
	softLinks := getSoftLinks(r.Header)

	return &uploadLocation{dirName, subDirectory, softLinks}
}

func convertPathToDirectory(path string) string {
	str := strings.Split(path, httpPathSeparator)
	return strings.Join(str, string(os.PathSeparator))
}

func convertDirectoryToPath(path string) string {
	str := strings.Split(path, string(os.PathSeparator))
	return strings.Join(str, httpPathSeparator)
}

func getSubDirectory(header http.Header) string {
	return header["X-Version"][0]
}

func getSoftLinks(header http.Header) []string {
	return strings.Split(header["X-Alias"][0], ",")
}

// TODO - change name to getDirectory from Path and add below changes
// Add documentation/Example about working
// its mandatory to have a / at the end
// if it does have a ., and does not end with a /, it will mean its a file location
func getFromPath(path string) (string, bool) {
	// strip out "/u/", then split by "/" to see the size
	newPath := strings.TrimLeft(path, UploadPrefix+httpPathSeparator)

	// remove blank strings
	data := strings.Split(newPath, httpPathSeparator)

	data = cleanStrings(data, "")
	return strings.Join(data, httpPathSeparator), (len(data) >= 2)
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

func getAction(ext string) ff.FileFormat {
	switch ext {
	// case "x-bzip2":
	// 	return bz2
	// case "x-gzip"
	//  return gzip
	// case "x-gtar":
	// 	return tar
	// tar.gz, .tgz, .tar.Z, .tar.bz2, .tbz2, .tar.lzma, .tlz
	// TODO - decide based on file format once saved
	// source info https://en.wikipedia.org/wiki/List_of_archive_formats
	case "zip":
		return ff.FileZip
		// .zip, .zipx
		// others to consider
		// rar
		// apk
		// jar
	case "json":
		return ff.FileJSON
	case "html":
		return ff.FileHTML
	case "plain":
		return ff.FileText
	default:
		return ff.FileText
	}
}

// // . Determine subdirectory to store based on other information
// switch action {
// case zip:
//   // save single file
//   // then unzip
//   // remove/save source file
// case html, text, json:
//   // save single file
//   // TODO - add for form submit type if json or encoding etc.,
// }
//

//UploadPrefix is required to get the upload prefix
const UploadPrefix = "/u"

// Prefix is the location at which the user can upload the files
func Prefix() string {
	return UploadPrefix + httpPathSeparator
}

type upload struct{}

func init() {
	u := &upload{}
	http.HandleFunc(Prefix(), u.delegate)
}

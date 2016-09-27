package upload

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

  Usecase discussion -
  1. A zip file is being uploaded with a dir/ | generate a version/filename or use the directory name
  2. A zip file is being uploaded with a file.zip | use the filename temporarily
  3. A zip file is being uploaded with a file.html | ignore filename and use our own file name to extract

  1. A html file being uploaded as file.html | use the filename to display the contents
  2. A html file being uploaded as file.zip | ignore/fail with error
  3. A html file being uploaded with a dir/ | use the default name like index.html or index.txt to upload

*/

func unableToDetermine(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Unable to determine the upload type. Internal Server Error!")
}

func (u *upload) delegate(w http.ResponseWriter, r *http.Request) {

	// Determine aliases types etc., Default is static.

	function := unableToDetermine
	// settings := &uploadSettings{undetermined, w, r}
	switch strings.Join(r.Header["X_Upload_Type"], "") {
	case "static", "Static", "onetime", "OneTime":
		function = Static // static files like zip or html without versioning (below code to SaveFile)
	case "versioned_files", "VersionedFiles":
		function = Versioned
	case "data_point", "data_points", "dataPoint", "dataPoints":
		function = DataPoints // data points like key/value one or multiple
		// data is usually appended to a file. like in parallel, but not necessarily large files
		// does not support/expect zips
	case "table", "Table":
		function = Table // json or CSV formats
		// should get a single file
	case "parallel", "Parallel":
		function = Parallel // multiple files to be put into the same location appended to each other
		// This needs to be passed to the SaveToFile parameter with a lockfile or something equivalent
		// only supports text files
	default:
		function = StaticNoOverride // in case of an existing directory in the place, we will throw an error

	}
	// Determine location to upload
	uploadLoc := generateUploadLocation(r)
	theaction := getAction(uploadLoc.extension)

	// move to directory and pop out
	dir := utils.MoveToFromDir("")
	fmt.Println(os.Getwd())
	os.Chdir(dir())
	fmt.Println(os.Getwd())
	defer os.Chdir(dir())

	// get main directory to save directory
	uploadLoc.generateDirectory()

	// read the data from the body based on type, save the file
	data, _ := ioutil.ReadAll(r.Body)
	utils.DebugHTTP(w, r)

	fmt.Println(uploadLoc)

	// save file in directory location
	_, err := utils.SaveToFile(uploadLoc.path(), data)
	if err != nil {
		fmt.Println(err)
	}
	output := action.Perform(theaction, uploadLoc.dirpath(), uploadLoc.filepath())

	// Aliases are to be made after the action is done.
	// This is to ensure failure of unzip or other misc actions do not point to a failed location
	subdir := utils.MoveToFromDir(uploadLoc.directory)
	os.Chdir(subdir())
	uploadLoc.makeAliases()
	os.Chdir(subdir())

	targetURL := config.ServerURL + display.Prefix() + convertDirectoryToPath(uploadLoc.dirpath())
	w.Header().Add("X-Generated-URL", targetURL)
	fmt.Fprintf(w, "\n"+output+"\n")

	_ = function //(w, r)
}

// Extract filename from header
func getFilename(r *http.Request, ext string, extfrom extFrom) (string, string) {
	var filename string

	// TODO - take care of extension as well or pass a separate header for it
	filename = r.Header.Get("X-File-Name")
	if filename != "" {
		return filename, ext
	}

	switch extfrom {
	case extURLPath:
		d := strings.Split(r.URL.Path, httpPathSeparator)
		lastPath := d[len(d)-1]
		filename = strings.Replace(lastPath, "."+ext, "", 1)

	case extContentType:
		// What to generate if extension comes from url is blank?
		// random or index?
		filename = "index"
	}

	return filename, ext
}

type extFrom int

const (
	extURLPath extFrom = iota
	extContentType
)

func getExt(r *http.Request) (string, extFrom) {
	d := strings.Split(r.URL.Path, httpPathSeparator)
	lastPath := d[len(d)-1]

	var extRegexp = regexp.MustCompile(`(tar\.gz|tar.bz2|gz|zip|bz2|txt|html|json|log)$`)
	var match = extRegexp.FindStringSubmatch(lastPath)
	if len(match) > 0 {
		return match[0], extURLPath
	}
	return extBasedOnContentType(r.Header["Content-Type"]), extContentType
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
	fmt.Println("-------------------")
	fmt.Println(dirName)
	if !ok {
		dirName = dirName + httpPathSeparator + utils.RandomString(config.RandomStringLength)
	}

	dirName = convertPathToDirectory(dirName)
	subDirectory := getSubDirectory(r.Header)
	softLinks := getSoftLinks(r.Header)

	// Later - Need to fork off here based on form vs upload type
	// If form data is submitted, we need to do action per file

	ext, extfrom := getExt(r)
	fileName, ext := getFilename(r, ext, extfrom)

	return &uploadLocation{dirName, subDirectory, softLinks, fileName, ext}
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
	return header.Get("X-Version")
}

func getSoftLinks(header http.Header) []string {
	return strings.Split(header.Get("X-Alias"), ",")
}

// TODO - change name to getDirectory from Path and add below changes
// Add documentation/Example about working
// its mandatory to have a / at the end
// if it does have a ., and does not end with a /, it will mean its a file location
func getFromPath(path string) (string, bool) {

	var ignoreLast = false
	if path[len(path)-1] != '/' {
		ignoreLast = true
	}

	// strip out "/u/", then split by "/" to see the size
	newPath := strings.TrimLeft(path, UploadPrefix+httpPathSeparator)

	// remove blank strings
	data := strings.Split(newPath, httpPathSeparator)
	if ignoreLast {
		data = data[:len(data)-1]
	}

	// TODO also, we need remove os delimiters invalid utf8 chars etc.,
	data = cleanStrings(data, "")
	return strings.Join(data, httpPathSeparator), (len(data) >= 2)
}

// cleanStrings skips array elements which are blank/nil
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

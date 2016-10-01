package upload

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	conf "../config"
	"../display"
	"../utils"
	"./action"
	ff "./fileformat/"
)

// "/" is being used at many locations for many operations. A Path separator makes sense since we have a file separator
// We need this if we want to run the server on cross platform later on.
const httpPathSeparator = "/"

type uploadType int

// Info ..
var Info = make(map[string]time.Time)

// Locker ..
var Locker = &sync.Mutex{}

const (
	undetermined uploadType = iota
	static                  // directory/file
	versioned               // directory/file
	dataPoint               // inside a directory. each data point is a file
	table                   // inside a directory, its a file
	parallel                // per file
)

// HeaderVersion is ..
const (
	HeaderVersion     = "X-Version" // HeaderVersion is used for creating documentation
	HeaderAlias       = "X-Alias"   // HeaderAlias used along side documentation for linking a version
	HeaderUploadType  = "X-Upload-Type"
	HeaderFileName    = "X-File-Name"
	HeaderAppend      = "X-Append" // used for Parallel/Concurrent writes
	HeaderContentType = "Content-Type"
	ResponseHeaderURL = "X-Generated-URL"
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
	switch strings.ToLower(r.Header.Get(HeaderUploadType)) {
	case "static", "onetime":
		function = Static // static files like zip or html without versioning (below code to SaveFile)
	case "versioned_file", "versionedfile":
		function = Versioned
	case "data_point", "datapoint":
		function = DataPoints // data points like key/value one or multiple
		// data is usually appended to a file. like in parallel, but not necessarily large files
		// does not support/expect zips
	case "table":
		function = Table // json or CSV formats
		// should get a single file
	case "parallel":
		function = Parallel // multiple files to be put into the same location appended to each other
		// This needs to be passed to the SaveToFile parameter with a lockfile or something equivalent
		// only supports text files
	default:
		function = StaticNoOverride // in case of an existing directory in the place, we will throw an error
	}

	var err error
	var uploadLoc *uploadLocation
	// read the data from the body based on type, save the file
	bodyData, _ := ioutil.ReadAll(r.Body)
	// utils.DebugHTTP(w, r)

	if datapoint := checkDataPointType(r.Header.Get(HeaderContentType)); datapoint != DataPointNone {
		dirName, ok := getFromPath(r.URL.Path)
		if !ok {
			dirName = dirName + httpPathSeparator + utils.RandomString(conf.C().Upload.DirectoryLength)
		}

		dirName = convertPathToDirectory(dirName)
		uploadLoc = &uploadLocation{dirName, "", []string{}, "", "datapoint"}

		dir := utils.MoveToFromDir("")
		os.Chdir(dir())
		defer os.Chdir(dir())

		uploadLoc.generateDirectory()

		// currentTime := time.Now().Unix()

		// This is a CSV
		r := csv.NewReader(strings.NewReader(string(bodyData)))
		records, err1 := r.ReadAll()
		var graphData [100]float64

		for i, record := range records {
			if len(record) >= 3 {
				filename := record[0]
				uploadLoc.filename = filename
				data := strings.Join(record[1:3], ",") + "\n"
				graphData[i], _ = strconv.ParseFloat(record[1], 64)
				utils.AppendToFile(uploadLoc.path(), []byte(data))
			}
		}
		if err1 != nil {
			log.Fatal(err1)
		}
		return
	}

	// Determine location to upload
	uploadLoc = generateUploadLocationForRawData(r)
	compressionType, archiveType := getAction(uploadLoc.extension)

	// move to directory and pop out
	dir := utils.MoveToFromDir("")
	os.Chdir(dir())
	defer os.Chdir(dir())

	uploadLoc.generateDirectory()
	// fmt.Println(uploadLoc)

	settings := &action.Settings{
		AppendMode:      false,
		CompressionType: compressionType,
		ArchiveType:     archiveType,
	}

	settings.AppendMode = isAppend(r.Header, settings)

	if settings.AppendMode {
		err = getLock(uploadLoc.path())
		if err != nil {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
		_, err = utils.AppendToFile(uploadLoc.path(), bodyData)
		releaseLock(uploadLoc.path())
	} else {
		_, err = utils.SaveToFile(uploadLoc.path(), bodyData)
	}

	// save file in directory location
	if err != nil {
		fmt.Println(err)
	}

	// _ = settings
	// output := "hello"
	action.Perform(settings, uploadLoc.dirpath(), uploadLoc.filepath())

	// Aliases are to be made after the action is done.
	// This is to ensure failure of unzip or other misc actions do not point to a failed location
	subdir := utils.MoveToFromDir(uploadLoc.directory)
	os.Chdir(subdir())
	uploadLoc.makeAliases()
	os.Chdir(subdir())

	targetURL := conf.C().Server.URL + display.Prefix() + convertDirectoryToPath(uploadLoc.dirpath())
	w.Header().Set(ResponseHeaderURL, targetURL)

	// fmt.Fprintf(w, "\n"+output+"\n")

	_ = function //(w, r)
}

func getLock(path string) error {
	i := 0
	for _, ok := Info[path]; ok; i++ {
		time.Sleep(5 * time.Millisecond)
		// wait for 10 times and timeout if lock is not released
		if i > 30 {
			break
		}
	}
	if i > 31 {
		return errors.New("Unable to unlock")
		// timeout with 504 and time taken to process the request
		// and the value of the path
	}
	Locker.Lock()
	// time.Sleep(500 * time.Millisecond) // test with this in case of doubt
	Info[path] = time.Now()
	Locker.Unlock()
	return nil

}
func releaseLock(path string) error {
	Locker.Lock()
	delete(Info, path)
	Locker.Unlock()
	return nil
}

// Extract filename from header
func getFilename(r *http.Request, ext string, extfrom extFrom) (string, string) {
	var filename string

	// TODO - take care of extension as well or pass a separate header for it
	filename = r.Header.Get(HeaderFileName)
	if filename != "" {
		ext1 := isWhiteListedRegexp(filename)
		if ext1 != "" {
			ext = ext1
			filename = strings.Replace(filename, "."+ext, "", 1)
		}
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

func isWhiteListedRegexp(str string) string {
	var extRegexp = regexp.MustCompile(`(tar\.gz|tar\.bz2|gz|zip|bz2|txt|html|json|log|png|jpe?g|bmp|svg|webp|md|markdown|toml|cfg|xml|xls)$`)
	var match = extRegexp.FindStringSubmatch(str)
	if len(match) > 0 {
		return match[0]
	}
	return ""
}

func getExt(r *http.Request) (string, extFrom) {
	d := strings.Split(r.URL.Path, httpPathSeparator)
	lastPath := d[len(d)-1]
	match := isWhiteListedRegexp(lastPath)
	if match != "" {
		return match, extURLPath
	}
	return extBasedOnContentType(r.Header.Get(HeaderContentType)), extContentType
}

// DataPoint is one of these types
// counting (incr/decr), value, gauges (can be increment from previous value), progress (means %)
type DataPoint int

// DataPointNone is None
const (
	DataPointNone = iota
	DataPointCounter
	DataPointValue
	DataPointProgress
)

func checkDataPointType(contentType string) DataPoint {
	str := strings.Split(contentType, httpPathSeparator)[1]
	var output DataPoint = DataPointNone
	switch str {
	// counter means we want to group by second/minute/hour
	case "vnd.datapoint+counter":
		output = DataPointCounter
		// value means we want to show it as is
	case "vnd.datapoint+value":
		output = DataPointValue
		// progress/percent may mean it wont be more than 100
		// TODO - percentage looks unnecessary. we can fit it in value itself
	case "vnd.datapoint+percentage":
		output = DataPointProgress
	}
	return output
}

func extBasedOnContentType(contentType string) string {
	str := strings.Split(contentType, httpPathSeparator)[1]
	output := str
	switch str {
	case "x-www-form-urlencoded":
		output = "txt"
	}
	return output
}

func generateUploadLocationForRawData(r *http.Request) *uploadLocation {
	/*
	  1.1. Get path
	  1.2. If path is not there, create a directory under level 2 like if uploaded to /u/docs, create /u/docs/{123456}/
	  1.3. If path is not there, create a directory under level 1 like if uploaded to /u, create /u/{123456}/
	*/

	dirName, ok := getFromPath(r.URL.Path)
	if !ok {
		dirName = dirName + httpPathSeparator + utils.RandomString(conf.C().Upload.DirectoryLength)
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
	return header.Get(HeaderVersion)
}

func getSoftLinks(header http.Header) []string {
	return strings.Split(header.Get(HeaderAlias), ",")
}

func isAppend(header http.Header, settings *action.Settings) bool {
	t := (strings.ToLower(header.Get(HeaderAppend)) == "true")
	if t && settings.CompressionType == ff.CompressionNone && settings.ArchiveType == ff.ArchiveNone {
		return true
	}
	return false
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

func getAction(ext string) (ff.CompressionFormat, ff.ArchiveFormat) {
	ext, compression := ff.FindCompressionFormat(ext)
	_, archive := ff.FindArchiveFormat(ext)
	return compression, archive
}

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

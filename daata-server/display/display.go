package display

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	conf "../config"
)

/*
 based on format
 1. Regular format is to display the data as-is
 2. Data Points format or JSON data format has key / value injected into a file
 3. Display appended file from multiple hosts
 4. Graphs / Dashboards
*/

// DisplayPrefix is a
const DisplayPrefix = "/d"

var dataPointRegexp = regexp.MustCompile(`\.datapoint$`)

func isDataPoint(str string) bool {
	return dataPointRegexp.MatchString(str)
}

func renderIfDataPoint(p string) ([]string, bool) {
	regularFlow := false
	f := config("directory") + p
	var files []string

	stat, err := os.Stat(f)
	// no such directory/file exists
	if err != nil {
		stat, err = os.Stat(f + ".datapoint")
		if err == nil {
			f += ".datapoint"
		}
	}

	if err == nil {
		if stat.IsDir() {
			// is a dir
			listFiles, _ := ioutil.ReadDir(f)
			for _, file := range listFiles {
				if isDataPoint(file.Name()) {
					files = append(files, f+file.Name())
				}
			}
			if len(files) == 0 {
				regularFlow = true
			}
		} else {
			// is a file is not a dir
			// display the data points
			if isDataPoint(f) {
				files = append(files, f)
			} else {
				regularFlow = true
			}
		}
	}
	return files, regularFlow
}

func openDir(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, DisplayPrefix)
	r.URL.Path = p
	files, regularFlow := renderIfDataPoint(p)

	if !regularFlow {
		for _, file := range files {
			Graph(w, r, file)
		}
	}

	if regularFlow {
		// check auth here
		// TODO - fix directory here from config
		http.FileServer(http.Dir(config("directory"))).ServeHTTP(w, r)
	}
}

func config(str string) string {
	if str == "directory" {
		return conf.C().Upload.Directory
	}
	return ""
}

// Prefix specifies the download/display location of a file
func Prefix() string {
	return DisplayPrefix + "/"
}

func init() {
	http.HandleFunc(Prefix(), openDir)
}

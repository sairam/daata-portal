package display

import (
	"net/http"
	"strings"

	"../config"
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

func openDir(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL.Path)
	if p := strings.TrimPrefix(r.URL.Path, DisplayPrefix); len(p) < len(r.URL.Path) {
		r.URL.Path = p
		// check auth here
		// TODO - fix directory here from config
		http.FileServer(http.Dir(config.DataDirectory)).ServeHTTP(w, r)
	} else {
		// fmt.Println(p)
		http.NotFound(w, r)
	}
}

// Prefix specifies the download/display location of a file
func Prefix() string {
	return DisplayPrefix + "/"
}

func init() {
	http.HandleFunc(Prefix(), openDir)
}

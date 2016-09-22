package main

import (
	"fmt"
	"net/http"
	"strings"
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
	if p := strings.TrimPrefix(r.URL.Path, DisplayPrefix); len(p) < len(r.URL.Path) {
		r.URL.Path = p
		// check auth here
		http.FileServer(http.Dir("../tmp")).ServeHTTP(w, r)
	} else {
		fmt.Println(p)
		http.NotFound(w, r)
	}
}

func init() {
	http.HandleFunc(DisplayPrefix+"/", openDir)
}

package main

import (
	"fmt"
	"net/http"
	"strings"
)

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

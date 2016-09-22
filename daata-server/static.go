package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

// StaticWebsiteInfo is of type string with a URL / Contact / Website
type StaticWebsiteInfo struct {
	URL     string
	Contact string
	Website string
}

// mapping "" to Helper probably maps to multiple urls including blank
// we should avoid this and find a cleaner way to map the home page.
var staticURLMappings = map[string]string{
	// "":        "usage", // index
	"help":    "usage",
	"about":   "about",
	"docs":    "docs",
	"contact": "contact",
	"pricing": "pricing",
	"404":     "404",
}

func page(name string, w io.Writer) {
	filename := staticURLMappings[name]
	if filename == "" {
		filename = staticURLMappings["help"]
	}
	filename += ".tmpl"
	t, err := template.New(filename).ParseFiles("tmpl/" + filename)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, StaticWebsiteInfo{URL: serverURL})
}

// StaticPage documentation
func StaticPage(w http.ResponseWriter, r *http.Request) {
	page(r.URL.Path, w)
}

func init() {
	for path := range staticURLMappings {
		http.HandleFunc("/"+path, StaticPage)
	}
}

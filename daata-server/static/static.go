package static

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"../config"
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
		fmt.Println("mapping did not match for " + name)
		filename = staticURLMappings["help"]
	}
	filename += ".tmpl"
	fmt.Println(config.StaticDirectory + "tmpl/" + filename)
	t, err := template.New(filename).ParseFiles(config.StaticDirectory + "tmpl/" + filename)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, StaticWebsiteInfo{URL: config.ServerURL})
}

// StaticPage documentation
func StaticPage(w http.ResponseWriter, r *http.Request) {
	path := "" + r.URL.Path
	path = strings.TrimLeft(path, "/")
	page(path, w)
}

func init() {
	for path := range staticURLMappings {
		http.HandleFunc("/"+path, StaticPage)
	}
}

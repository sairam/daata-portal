package static

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	conf "../config"
)

// WebsiteInfo is of type string with a URL / Contact / Website
type WebsiteInfo struct {
	URL     string
	Contact string
	Website string
}

// mapping "" to Helper probably maps to multiple urls including blank
// we should avoid this and find a cleaner way to map the home page.
var urlMappings = map[string]string{
	// "":        "usage", // index
	"help":    "usage",
	"about":   "about",
	"docs":    "docs",
	"contact": "contact",
	"pricing": "pricing",
	"404":     "404",
}

func page(name string, w io.Writer) {
	filename := urlMappings[name]
	if filename == "" {
		fmt.Println("mapping did not match for " + name)
		filename = urlMappings["help"]
	}
	filename += ".tmpl"
	// fmt.Println(config.StaticDirectory + "tmpl/" + filename)
	t, err := template.New(filename).ParseFiles(conf.C().Directories.Static + "tmpl/" + filename)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.Execute(w, WebsiteInfo{URL: conf.C().Server.URL})
}

// Page documentation
func Page(w http.ResponseWriter, r *http.Request) {
	path := "" + r.URL.Path
	path = strings.TrimLeft(path, "/")
	w.Header().Set("Content-Type", "text/html;utf8")
	page(path, w)
}

func init() {
	for path := range urlMappings {
		http.HandleFunc("/"+path, Page)
	}
}

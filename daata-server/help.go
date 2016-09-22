package main

import (
	"html/template"
	"io"
)

// Data is of type string with a URL
type Data struct {
	URL string
}

func help(w io.Writer) {
	filename := "usage.tmpl"
	t, _ := template.New(filename).ParseFiles("tmpl/" + filename)
	t.Execute(w, Data{URL: serverURL})
}

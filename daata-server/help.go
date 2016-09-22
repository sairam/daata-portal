package main

import (
	"html/template"
	"net/http"
)

// Data is of type string with a URL
type Data struct {
	URL string
}

func help(w http.ResponseWriter, _ *http.Request) {

	t, _ := template.New("foo").Parse(`
curl -i -X POST {{.URL}}/upload -H "Content-Type: application/zip" --data-binary "@data.zip"
curl -i -X POST {{.URL}}/upload -H "Content-Type: application/json" --data-binary "@freshmenu.json"

Examples:
{{.URL}}/d/wwpbbi
{{.URL}}/d/iwhspu
`)
	t.Execute(w, Data{URL: serverURL})
}

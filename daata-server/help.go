package main

import (
	"html/template"
	"io"
)

// Data is of type string with a URL
type Data struct {
	URL string
}

const helpInfo = `
curl -i -X POST {{.URL}}/upload -H "Content-Type: application/zip" --data-binary "@data.zip"
curl -i -X POST {{.URL}}/upload -H "Content-Type: application/json" --data-binary "@freshmenu.json"

Examples:
{{.URL}}/d/wwpbbi
{{.URL}}/d/iwhspu
`

func help(w io.Writer) {
	t, _ := template.New("").Parse(helpInfo)
	t.Execute(w, Data{URL: serverURL})
}

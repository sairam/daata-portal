package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const randomStringLength = 6
const port = ":8001"

// Fileformat should be used for understand file extension
// fdlsak
type FileFormat int

const (
	text FileFormat = iota + 1
	json
	html
	zip
)

func main() {
	fmt.Printf("Hello Server is on localhost%s\n", port)
	// action := getAction("zip")
	// performAction(action, "../tmp/test12", "test12.zip")
	// fmt.Println(randomString(randomStringLength))

	http.HandleFunc("/d/", openDir)
	http.HandleFunc("/upload", saveFile)
	http.HandleFunc("/", help)
	http.ListenAndServe(port, nil)
}

func openDir(w http.ResponseWriter, r *http.Request) {
	prefix := "/d"
	if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
		r.URL.Path = p
		// check auth here
		http.FileServer(http.Dir("../tmp")).ServeHTTP(w, r)
	} else {
		fmt.Println(p)
		http.NotFound(w, r)
	}
}

//
// func showFile(w http.ResponseWriter, r *http.Request) {
// 	fileName := r.URL.String()
// 	data := strings.TrimLeft(fileName, "/d/")
// 	openIndex(fileName)
// 	fmt.Fprintf(w, "%s\n", data)
// }
//
// func openIndex(file string) {
//
// }

const serverURL = "https://c33a453f.ap.ngrok.io"

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

// "http://localhost" + port
func saveFile(w http.ResponseWriter, r *http.Request) {
	// 0. generate random id
	dirName := randomString(randomStringLength)
	url := serverURL + "/d/" + dirName
	// 1. read contents
	data, _ := ioutil.ReadAll(r.Body)
	debug(w, r)

	// 2. save file
	extension := strings.Split(r.Header["Content-Type"][0], "/")[1]
	directory, fileName := saveToFile(dirName, extension, data)

	// 3. determine file type
	action := getAction(extension)

	// 4. perform action of unzip or nothing
	output := performAction(action, directory, fileName)
	fmt.Fprintf(w, "\n"+output+"\n")

	// 5. send back url based on random id
	fmt.Fprintf(w, "\n"+url+"\n")
}

func performAction(format FileFormat, dir, file string) string {
	currentDirectory, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(currentDirectory)
	if format == zip {
		cmd := []string{"/usr/bin/unzip", file}
		out, err := exec.Command(cmd[0], cmd[1]).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Output is \n%s\n", out)
		return string(out)
	}
	return ""
}

func getAction(ext string) FileFormat {
	if ext == "zip" {
		return zip
	}
	return text
}

func saveToFile(name, extension string, data []byte) (string, string) {
	dir := "../tmp/" + name
	fileName := name + "." + extension
	os.Mkdir(dir, 0700)
	ioutil.WriteFile(dir+"/"+fileName, data, 0600)
	fmt.Println("Created file")
	return dir, fileName
}

func debug(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Hos = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

func randomString(length int) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyz" // "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

/*
Later
Access restriction
If unzip file does not contain index.html, generate one with a tree.
*/

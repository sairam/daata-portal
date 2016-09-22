package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const randomStringLength = 6
const port = ":8001"
const serverURL = "https://21ae9584.ap.ngrok.io"

//RedirectPrefix is required to generate redirect links
const RedirectPrefix = "/r"

func main() {
	fmt.Printf("Hello Server is on localhost%s\n", port)
	// action := getAction("zip")
	// performAction(action, "../tmp/test12", "test12.zip")
	// fmt.Println(randomString(randomStringLength))

	http.HandleFunc(RedirectPrefix+"/", Redirect)
	http.HandleFunc("/d/", openDir)
	http.HandleFunc("/upload", saveFile)
	http.HandleFunc("/", Help)
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

// Upload data
func Upload(w http.ResponseWriter, r *http.Request) {
	saveFile(w, r)
}

// Help data
func Help(w http.ResponseWriter, _ *http.Request) {
	help(w)
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

/*
Later
Access restriction
If unzip file does not contain index.html, generate one with a tree.
*/

package utils

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"../config/"
)

// MoveToFromDir ..
// generates a function which returns target dir and then current dir
func MoveToFromDir(str string) func() string {
	i := 0
	pwd, _ := os.Getwd()
	dataDir := config.DataDirectory + "/" + str
	return func() string {
		if i == 0 {
			i++
			return dataDir
		}
		return pwd
	}
}

// DebugHTTP ..
func DebugHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// SaveToFile ..
func SaveToFile(filename string, data []byte) (string, error) {
	// filepath := strings.Join([]string{parentDir, filename}, string(os.PathSeparator))
	err := ioutil.WriteFile(filename, data, 0600)
	// wd, _ := os.Getwd()
	// fmt.Printf("Created file at %s, %s\n", filename, wd)
	return filename, err
}

// RandomString ..
func RandomString(length int) string {
	const alphanum = "abcdefghijklmnopqrstuvwxyz" // "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var bytes = make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

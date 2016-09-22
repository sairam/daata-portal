package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

func saveToFile(name, extension string, data []byte) (string, string) {
	dir := "../tmp/" + name
	fileName := name + "." + extension
	os.Mkdir(dir, 0700)
	ioutil.WriteFile(dir+"/"+fileName, data, 0600)
	fmt.Println("Created file")
	return dir, fileName
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

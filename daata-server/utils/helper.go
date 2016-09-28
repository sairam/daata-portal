package utils

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"../config/"
)

const b62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// copied example code from http://chartd.co/
func EncodeGraphData(data []float64, min, max float64) string {
	r := math.Dim(max, min)
	bs := make([]byte, len(data))
	if r == 0 {
		for i := 0; i < len(data); i++ {
			bs[i] = b62[0]
		}
		return string(bs)
	}
	enclen := float64(len(b62) - 1)
	for i, y := range data {
		index := int(enclen * (y - min) / r)
		if index >= 0 && index < len(b62) {
			bs[i] = b62[index]
		} else {
			bs[i] = b62[0]
		}
	}
	return string(bs)

}

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

func saveToTempLocation(data []byte) (string, string, error) {
	dir, err := ioutil.TempDir("", "upload")
	file := RandomString(8)
	err = ioutil.WriteFile(dir+file, data, 0600)
	return dir, file, err
}

// AppendToFile ..
func AppendToFile(filename string, data []byte) (string, error) {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return filename, err
	}

	defer f.Close()

	_, err = f.Write(data)
	return filename, err
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

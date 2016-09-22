package main

import (
	"fmt"
	"net/http"
)

const randomStringLength = 6
const port = ":8001"
const serverURL = "https://21ae9584.ap.ngrok.io"
const dataDirectory = "../data/"
const maxUploadParamsLimit = 5000

// This is to be used in all upload forms

func main() {
	fmt.Printf("Hello Server is on localhost%s\n", port)

	//  / - index
	//  /help - help
	//  /d/ - display uploaded data
	//  /g/ - graphs uses data of a particular format to make simple line/bar graphs
	//  /u/ - upload
	//  /r/ - redirects

	http.ListenAndServe(port, nil)
}

package main

import (
	"fmt"
	"net/http"
)

const randomStringLength = 6
const port = ":8001"
const serverURL = "https://21ae9584.ap.ngrok.io"
const dataDirectory = "../data/"

func main() {
	fmt.Printf("Hello Server is on localhost%s\n", port)

	/*
	   / - index
	   /help - help
	   /d/ - /display
	   /u/ - upload
	   /r/ - redirects
	*/

	http.ListenAndServe(port, nil)
}

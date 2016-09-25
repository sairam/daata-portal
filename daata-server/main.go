package main

import (
	"fmt"
	"net/http"

	"./config"
	_ "./display"
	_ "./redirect"
	_ "./static"
	_ "./upload"
	_ "./utils"
)

// This is to be used in all upload forms

func main() {
	fmt.Printf("Hello Server is on localhost%s\n", config.Port)

	//  / - index
	//  /help - help
	//  /d/ - display uploaded data
	//  /g/ - graphs uses data of a particular format to make simple line/bar graphs
	//  /u/ - upload
	//  /r/ - redirects

	http.ListenAndServe(config.Port, nil)
}

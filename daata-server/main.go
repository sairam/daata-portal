package main

import (
	"fmt"
	"log"
	"net/http"

	conf "./config"
	_ "./display"
	_ "./redirect"
	_ "./static"
	_ "./upload"
	// _ "./utils"

	_ "net/http/pprof"
)

// This is to be used in all upload forms

func main() {
	port := fmt.Sprintf("%d", conf.C().Server.Port)
	fmt.Printf("Server is being served on http://localhost:%s\n", port)

	//  / - index
	//  /help - help
	//  /d/ - display uploaded data
	//  /g/ - graphs uses data of a particular format to make simple line/bar graphs
	//  /u/ - upload
	//  /r/ - redirects

	log.Println(http.ListenAndServe(":"+port, nil))
}

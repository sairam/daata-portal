package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

//FileFormat is a should be used for understand file extension
type FileFormat int

const (
	text FileFormat = iota + 1
	json
	html
	zip
)

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

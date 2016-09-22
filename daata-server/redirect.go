package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// TODO use these structs
type urlShortner struct {
	shortURL string
	longURL  string
}

type urlShortnerForm struct {
	urlShortner
	Override bool
}

/* Query Data Store */

// returns true if the file exists
// returns false if the file does not exist
func exists(shortURL string) bool {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	if _, err := os.Stat(shortURL); os.IsNotExist(err) {
		return false
	}
	return true
}

func read(shortURL string) (string, error) {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	data, err := ioutil.ReadFile(shortURL)
	if err != nil {
		return "", err
	}

	dataStr := string(data)
	longURL := strings.Split(dataStr, "\n")[0]

	return longURL, nil
}

func moveToDir() func() string {
	// this is the file system prefix
	return moveToFromDir(RedirectPrefix + "/")
}

func insert(shortURL, longURL string) error {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	file, err := os.OpenFile(shortURL, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(longURL + "\n")
	return nil
}

func update(shortURL, longURL string) error {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	file, err := os.OpenFile(shortURL, os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString(longURL + "\n")
	return nil
}

/* Service Layer */
// shortURL without protocol, just a string of [a-zA-Z0-9-/] are allowed
// longURL can be anything including http, https, itunes urls or any other.
// for now limit to http, https for regular users.
// leading and trailing slashes will be removed.

// insert query, don't update
func insertshortURL(shortURL, longURL string) error {
	return makeEntryEvenIfExists(shortURL, longURL, false)
}

// upsert query
func upsertshortURL(shortURL, longURL string) error {
	return makeEntryEvenIfExists(shortURL, longURL, true)
}

func noOp(_, _ string) error {
	return errors.New("could not process")
}

func makeEntryEvenIfExists(shortURL, longURL string, override bool) error {
	function := noOp
	if exists(shortURL) {
		if override {
			function = update
		}
	} else {
		function = insert
	}
	return function(shortURL, longURL)
}

// CreateOrUpdateURL is the main method to add a new redirect
func CreateOrUpdateURL(shortURL, longURL string, update bool) (string, error) {
	if shortURL == "" {
		shortURL = randomString(6)
	}
	var err error
	if update {
		err = upsertshortURL(shortURL, longURL)
	} else {
		err = insertshortURL(shortURL, longURL)
	}
	return shortURL, err
}

// add caching if required in service layer
func findRedirectURL(shortURL string) (string, error) {
	if shortURL != "" && exists(shortURL) {
		return read(shortURL)
	}
	return "", errors.New("no such url exists")
}

func stripPrefix(path string) string {
	if p := strings.TrimPrefix(path, RedirectPrefix); len(p) < len(path) {
		if p[0] == '/' {
			return p[1:]
		}
		return p
	}
	return ""
}

// TODO - check how to make this cleaner
func validate(shortURL, longURL string) error {
	var err error

	err = validateShortURL(shortURL)
	if err != nil {
		return err
	}

	err = validateBlankURL(longURL)
	if err != nil {
		return err
	}

	valid, err := validateLongURL(longURL)
	if err != nil {
		return err
	}

	if !valid {
		err = validateRelativePath(longURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateShortURL(_ string) error {
	// ensure there are no spaces, dots or any such
	// whitelist with unicode chars.
	// TODO - make a demo with emojicons
	// validate if its a valid file system path
	return nil
}

func validateBlankURL(str string) error {
	if str == "" {
		return errors.New("long_url is blank")
	}
	return nil
}

func validateRelativePath(str string) error {
	if str[0] != '/' {
		return errors.New("url does not start with '/'")
	}
	// ensure does not have script tag
	return nil
}

func validateLongURL(_ string) (bool, error) {
	// TODO check if url starts with http or https
	return true, nil
}

func stripIfRedirect(path string) (string, bool) {
	redirect := true
	length := len(path) - 1
	if path[length] == '+' {
		redirect = false
		path = path[:length]
	}
	return path, redirect
}
func parseOverride(str string) bool {
	return (str == "true")
}

// Redirect is the main method which takes care of redirecting
// TODO Check Auth
func Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := stripPrefix(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		path, redirect := stripIfRedirect(shortURL)
		url, err := findRedirectURL(path)

		if err != nil || url == "" {
			http.NotFound(w, r)
		} else {
			if redirect {
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			} else {
				fmt.Fprintf(w, url)
			}
		}
	case http.MethodPost:
		// err := r.ParseForm()
		err := r.ParseMultipartForm(maxUploadParamsLimit)
		if err != nil {
			// TODO - generate form based on Content-Type
			// application/x-www-form-urlencoded
			// multipart/form-data
			// application/json
			// request Content-Type isn't multipart/form-data if r.ParseForm()
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		shortURL, longURL := r.Form.Get("short_url"), r.Form.Get("long_url")
		override := parseOverride(r.Form.Get("override"))
		url, err := CreateOrUpdateURL(shortURL, longURL, override)

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			fmt.Fprintf(w, url)
		}

	default:
		http.NotFound(w, r)
	}
}

//RedirectPrefix is required to generate redirect links
const RedirectPrefix = "/r"

func init() {
	http.HandleFunc(RedirectPrefix+"/", Redirect)
}

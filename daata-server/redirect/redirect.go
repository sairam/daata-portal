package redirect

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"../config/"
	"../utils/"
)

type persistInterface interface {
	exists() bool
	read()
	insert() error
	update() error
}

// TODO use these structs
type urlShortner struct {
	shortURL string
	longURL  string
	override bool
}

/* Query Data Store */

// returns true if the file exists
// returns false if the file does not exist
func (u *urlShortner) exists() bool {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	if _, err := os.Stat(u.shortURL); os.IsNotExist(err) {
		return false
	}
	return true
}

func (u *urlShortner) read() error {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	data, err := ioutil.ReadFile(u.shortURL)
	if err != nil {
		return err
	}

	dataStr := string(data)
	u.longURL = strings.Split(dataStr, "\n")[0]

	return nil
}

func (u *urlShortner) insert() error {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	file, err := os.OpenFile(u.shortURL, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(u.longURL + "\n")
	return nil
}

func (u *urlShortner) update() error {
	dir := moveToDir()
	os.Chdir(dir())
	defer os.Chdir(dir())

	file, err := os.OpenFile(u.shortURL, os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString(u.longURL + "\n")
	return nil
}

func moveToDir() func() string {
	// this is the file system prefix
	return utils.MoveToFromDir(RedirectPrefix + "/")
}

type serviceInterface interface {
	Validate() error
	CreateOrUpdate() error
}

/* Service Layer */
// shortURL without protocol, just a string of [a-zA-Z0-9-/] are allowed
// longURL can be anything including http, https, itunes urls or any other.
// for now limit to http, https for regular users.
// leading and trailing slashes will be removed.

func (u *urlShortner) noOp() error {
	return errors.New("could not process")
}

// CreateOrUpdateURL is the main method to add a new redirect
func (u *urlShortner) CreateOrUpdate() error {
	if u.shortURL == "" {
		u.shortURL = utils.RandomString(6)
	}
	// insert on false
	// update on true

	if u.exists() {
		if u.override {
			return u.update()
		}
	} else {
		return u.insert()
	}
	return u.noOp()
}

// add caching if required in service layer
func (u *urlShortner) Find() error {
	if u.shortURL != "" && u.exists() {
		return u.read()
	}
	return errors.New("no such url exists")
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
func (u *urlShortner) Validate() error {
	var err error

	err = validateShortURL(u.shortURL)
	if err != nil {
		return err
	}

	err = validateBlankURL(u.longURL)
	if err != nil {
		return err
	}

	valid, err := validateLongURL(u.longURL)
	if err != nil {
		return err
	}

	if !valid {
		err = validateRelativePath(u.longURL)
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
	// TODO - look at StripPrefix handler
	shortURL := stripPrefix(r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		path, redirect := stripIfRedirect(shortURL)
		u := &urlShortner{shortURL: path}
		err := u.Find()
		url := u.longURL

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
		err := r.ParseMultipartForm(config.MaxUploadParamsLimit)
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
		u := &urlShortner{shortURL: shortURL, longURL: longURL, override: override}
		err = u.CreateOrUpdate()

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			fmt.Fprintf(w, u.shortURL)
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

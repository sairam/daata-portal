package redirect

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/afero"

	conf "../config/"
	"../utils/"
)

// model
type urlShortner struct {
	shortURL string
	longURL  string
	override bool
}

// This interface for using the service
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
	p := strings.TrimPrefix(path, RedirectPrefix)
	if p[0] == '/' {
		return p[1:]
	}
	return p
}

func stripIfRedirect(path string) (string, bool) {
	redirect := true
	length := len(path) - 1
	if len(path)-1 < 0 {
		return path, false
	}
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
		if path == "" {
			// display 'new' page
			http.NotFound(w, r)
		} else {
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
		}
	case http.MethodPost:
		// err := r.ParseForm()
		err := r.ParseMultipartForm(int64(conf.C().Upload.SizeLimit))
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

var appFs afero.Fs
var fsutil *afero.Afero

func init() {
	// "r" is the directory

	appFs = afero.NewBasePathFs(afero.NewOsFs(), conf.C().Redirect.Directory)
	fsutil = &afero.Afero{Fs: appFs}

	http.HandleFunc(RedirectPrefix+"/", Redirect)
}

package redirect

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/afero"

	conf "../config/"
)

// model
type urlShortner struct {
	shortURL string
	longURL  string
	override bool
}

//RedirectPrefix is required to generate redirect links
const (
	RedirectPrefix             = "/r"
	AutoGenerateShortURLLength = 6
	ParamShortURL              = "short_url"
	ParamLongURL               = "long_url"
	ParamOverride              = "override"
)

var appFs afero.Fs
var fsutil *afero.Afero

func fsSettings(val string) {
	switch val {
	case "memory":
		appFs = afero.NewMemMapFs()
	case "readFS":
		appFs = afero.NewReadOnlyFs(afero.NewMemMapFs())
	case "writeFS":
		fallthrough
	default:
		appFs = afero.NewBasePathFs(afero.NewOsFs(), conf.C().Redirect.Directory)
	}
	fsutil = &afero.Afero{Fs: appFs}
}

func init() {
	fsSettings("writeFS")
	http.HandleFunc(RedirectPrefix+"/", Redirect)
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
			var (
				u                        = &urlShortner{shortURL: path}
				service serviceInterface = u
			)
			err := service.Find()
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
		var (
			shortURL                  = r.Form.Get(ParamShortURL)
			longURL                   = r.Form.Get(ParamLongURL)
			override                  = parseOverride(r.Form.Get(ParamOverride))
			u                         = &urlShortner{shortURL, longURL, override}
			service  serviceInterface = u
		)
		err = service.CreateOrUpdate()

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		} else {
			fmt.Fprintf(w, u.shortURL)
		}

	default:
		http.Error(w, "", http.StatusUnprocessableEntity)
		// http.NotFound(w, r)
	}
}

// Helper methods
func stripPrefix(path string) string {
	p := strings.TrimPrefix(path, RedirectPrefix)
	if p != "" && p[0] == '/' {
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
	return (strings.ToLower(str) == "true")
}

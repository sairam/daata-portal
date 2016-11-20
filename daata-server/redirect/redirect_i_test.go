package redirect

// This file contains URL endpoint test
import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"../utils"

	"github.com/spf13/afero"
)

const redirectURL = "r/"

func TestMain(m *testing.M) {
	appFs = afero.NewMemMapFs()
	fsutil = &afero.Afero{Fs: appFs}
	// defer cleanup()
	// TODO setup config

	os.Exit(m.Run())
}

func TestSubmitActualForm(t *testing.T) {
	// req := httptest.NewRequest("POST", "http://example.com/foo", nil)
	// add headers

	// w := httptest.NewRecorder()
	// Redirect(w, req)
	//
	// fmt.Printf("%d - %s", w.Code, w.Body.String())
}

func submitForm(url string, kv url.Values, file string) (response *http.Response, err error) {
	// Prepare a form that you will submit to that URL.
	var fw io.Writer

	response = &http.Response{}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = file
	// Add your image file
	// if file != "" {
	// 	f, err := os.Open(file)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer f.Close()
	//
	// 	fw, err = w.CreateFormFile("image", file)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if _, err = io.Copy(fw, f); err != nil {
	// 		return err
	// 	}
	// }

	// Add the other fields
	for k, v := range kv {
		if fw, err = w.CreateFormField(k); err != nil {
			return
		}
		if _, err = fw.Write([]byte(v[0])); err != nil {
			return
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	return res, nil

	// fmt.Printf("status: %s\n", res.Status)
	//
	// // Check the response
	// if res.StatusCode != http.StatusOK {
	// 	err = fmt.Errorf("bad status: %s", res.Status)
	// }
	// return
}
func fetchResponse(url string) (response *http.Response, err error) {
	return http.Get(url)
}

func TestUnKnownPath(t *testing.T) {
	shortURL := utils.RandomString(4)

	ts := httptest.NewServer(http.HandlerFunc(Redirect))
	urlEp := ts.URL + "/"
	defer ts.Close()
	response, _ := fetchResponse(urlEp + shortURL + "+")
	if response.StatusCode != 404 {
		t.Errorf("Magically found a unicorn! ")
	}

}

func TestUnsupportedPUTCall(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(Redirect))
	urlEp := ts.URL + "/"
	defer ts.Close()

	data1 := map[string]string{
		"short_url": "shortURL",
		"long_url":  "https://www.example.com",
		"override":  "true",
	}

	kv := url.Values{}
	for k, t := range data1 {
		kv.Set(k, t)
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	var err error

	for k, v := range kv {
		if fw, err = w.CreateFormField(k); err != nil {
			return
		}
		if _, err = fw.Write([]byte(v[0])); err != nil {
			return
		}
	}

	req, err := http.NewRequest("PUT", ts.URL+"/r/", &b)

	response, err := submitForm(urlEp, kv, "")
	ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		t.Errorf("error submitting form %s", err)
	}

	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, _ := client.Do(req)
	if res.StatusCode != 422 {
		t.Errorf("PUT should not be processable")
	}
}

func TestWithOverrideTrue(t *testing.T) {
	shortURL := utils.RandomString(4)
	data := map[string]string{
		"short_url": shortURL,
		"long_url":  "https://www.example.com",
		"override":  "true",
	}
	response := makeTheCall(redirectURL, data, t)

	if response != data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}
}

func TestWithoutOverride(t *testing.T) {
	shortURL := utils.RandomString(4)
	data := map[string]string{
		"short_url": shortURL,
		"long_url":  "https://www.sairam.com",
	}
	response := makeTheCall(redirectURL, data, t)

	if response != data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}
}

func TestWithOverrideFalse(t *testing.T) {
	shortURL := utils.RandomString(4)

	data := map[string]string{
		"short_url": shortURL,
		"long_url":  "https://www.example.com",
		"override":  "false",
	}

	response := makeTheCall(redirectURL, data, t)

	if response != data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}

	data = map[string]string{
		"short_url": shortURL,
		"long_url":  "https://www.example.in",
		"override":  "false",
	}
	response = makeTheCall(redirectURL, data, t)

	if response == data["long_url"] {
		t.Errorf("Problem. was able to override %s | %s", response, data["long_url"])
	}
}

func makeTheCall(urlA string, data map[string]string, t *testing.T) string {

	ts := httptest.NewServer(http.HandlerFunc(Redirect))
	urlEp := ts.URL + "/" + urlA
	defer ts.Close()

	v := url.Values{}
	for k, t := range data {
		v.Set(k, t)
	}
	response, err := submitForm(urlEp, v, "")
	ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		t.Errorf("error submitting form %s", err)
	}

	response, err = fetchResponse(urlEp + data["short_url"] + "+")

	if err != nil {
		t.Errorf("error fetching response is %s", err)
	}

	ding, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", ding)
}

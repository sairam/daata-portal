package redirect_i_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
)

const host = "http://localhost:8001/"
const redirectURL = "r/"

func init() {

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
func TestWithOverride(t *testing.T) {
	data := map[string]string{
		"short_url": "test",
		"long_url":  "https://www.example.com",
		"override":  "true",
	}
	response := actualTest(host+redirectURL, data, t)

	if response != data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}
}

func TestWithoutOverride(t *testing.T) {
	data := map[string]string{
		"short_url": "test",
		"long_url":  "https://www.sairam.com",
	}
	response := actualTest(host+redirectURL, data, t)

	if response == data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}
}

func TestWithFalseOverride(t *testing.T) {
	data := map[string]string{
		"short_url": "test",
		"long_url":  "https://www.data.com",
		"override":  "false",
	}
	response := actualTest(host+redirectURL, data, t)

	if response == data["long_url"] {
		t.Errorf("Problem. unable to set %s | %s", response, data["long_url"])
	}
}

func actualTest(urlA string, data map[string]string, t *testing.T) string {

	v := url.Values{}
	for k, t := range data {
		v.Set(k, t)
	}
	response, err := submitForm(urlA, v, "")
	response, err = fetchResponse(urlA + data["short_url"] + "+")

	if err != nil {
		t.Errorf("ererrrr is %s", err)
	}
	ding, _ := ioutil.ReadAll(response.Body)
	return string(ding)
}

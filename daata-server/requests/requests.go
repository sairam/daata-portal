package requests

import (
	"io"
	"net/http"
	"strings"
)

func displayData(w io.Writer, url string) {

}

// All misc requests go into this package
// 1. This includes handling incoming webhooks
// 2. Proxying/Forwarding webhooks
// 3. Receving dumps of requests
// 4. Forwarding or replaying the requests to other locations
// 2,4 have replaying in common
// 1,2 have requesting from external party in common
// (1,2), (3,4) have the format in common. The type can be a webhook or a request.

// Dumper is the main method which recieves requests
// The urls can be the endpoints which can be used to share with other teams to test endpoint checks for both request/response playbacks.
// This can be considered a similar tool to Postman.
// Usecase 1:
// A new API is exposed by Team B which Team A will be calling.
// Both services are hosted on the cloud. Team A will be sending their request to an endpoint say, /r/abcdef
// This will dump the requests. Now you can forward the request to Team B with the original payload from Team A including headers etc.,
// They can change the payload and we can diff across changes for user to easily identify changes needed to make in their code
// This is useful when there are no libraries exposed by Team B for Team A.
//
// Internal Workings
// 1. Each new request is considered as a new version.
// 2. Saving a new request will override the existing id while keeping the revision history (how to maintain revision history?)
// 3. When a 'New Request' is made, we create an endpoint. Data should be sent to that endpoint `/w/abcdef`.
// 4. Any POST, PUT, OPTIONS, GET calls made get logged with the request format including headers at `/w/abcdef/1475545366120`
// 5. Request IDs are made in milliseconds (unixtimestamp + millisecond). All edits will have the format 1475545366120-1, 1475545366120-2 for all edits made from the UI
// 6. Where this fails, something that comes within the same millisecond
func Dumper(w http.ResponseWriter, r *http.Request) {
	url = strings.TrimPrefix(r.URL.Path, WebhookPrefix)
	// Save the dump of headers, request body etc.,
	request := saveRequest(w, url, r)
	if targetURL, sync := getProxy(url); targetURL != "" {
		Proxy(w, r, targetURL, sync)
	}
	// We are done with the request
}

// Proxy is the actual request forwarder
// Proxy sync vs async
// There are 3 ways to configure a proxy (this is the priority order)
// 1. sending the url in the request in the header
// 2. sending the url in the request after the endpoint
// 3. configuring the proxy/forwarding in the web UI which is saved as .proxy file. Note: this is per endpoint which will call /w/abcdef/proxy
//    Any translations/changes of data will occur here
// (3) is where service requests like hidden proxy will fit in. Useful for hiding information.
// This is where little dynamic part comes in to match Tokens etc., which can be part of the YAML/toml configuration in the .proxy file
// A proxy can be set and called with /w/endpoint/http://example.com/path1/path2?arg=value
// the path after our endpoint will act like a proxy
func Proxy(w http.ResponseWriter, r *http.Request, targetURL string, sync bool) {
	// Based on sync/async, we should be able to spawn off a go routine
	if sync {
		ProxySync(w, r, targetURL)
	} else {
		// if async, send a 200 to the original request
		go Proxy(request, targetURL)
	}
}

func proxySync(w http.ResponseWriter, r *http.Request, targetURL string) {

	return
}

func proxyASync(r *http.Request, targetURL string) {
	return
}

// getProxy is service layer function which returns the url to send the request to
func getProxy(url string, sync bool) string {
	// checks if the request has requests in any of the 3 formats
	// returns only the url.
	// NOTE: a url cannot be of the format /w/... This is to avoid a redirect loop hell.
	// Having any other path is allowed. This can be used to test other metrics or services provided.
	return "", true
}

// Relay is responsible to send the request based on the data saved in dumper
// and send to the target endpoint.
// When this is automated with a script or something, it will be considered a proxy
func Relay(w http.ResponseWriter, r *http.Request) {

}

func DumperView(w http.ResponseWriter, r *http.Request) {
	url = strings.TrimPrefix(r.URL.Path, viewPath)
	switch r.Method {
	case http.MethodGet:
		displayData(w, url)
	case http.MethodPost:
		return saveRequest(w, url, r)
	}
}

// WebhookPrefix is used for catching and replaying requests.
const WebhookPrefix = "/w"
const viewPath = WebhookPrefix + "/view/"

func init() {
	http.HandleFunc(viewPath, DumperView)
	http.HandleFunc(WebhookPrefix+"/", Dumper)
}

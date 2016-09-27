package upload

import "net/http"

// Parallel is
func Parallel(w http.ResponseWriter, r *http.Request) error {
	return nil

}

/*
  uploading static files takes only static files hosted under a directory.
  If a directory is not specified, one will be allocated to it.

  Future callers should call that directory

  Static:
  curl -i -X POST -H "Content-Type: application/zip" --data-binary "@data.zip" https://my.daata.xyz/u/docs/spokes-platform

*/

// Static ..
func Static(w http.ResponseWriter, r *http.Request) error {

	return nil
}

// DataPoints is
func DataPoints(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// Table is
func Table(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Versioned Documents
// Header - Alias
// Header - Version
// Content-Type - zip
// PATH to upload
// Type - POST
// File-binary

/*

curl \
  -F "userid=1" \
  -F "filecomment=This is an image file" \
  -F "image=@/home/user1/Desktop/test.jpg" \
  localhost/uploader.php

curl -X POST \
  -H 'version: "2.1.3"' \
  -H 'alias: ["release-20160707", "master", "stable"]' \
  -H 'authorization: "abcdefghijklmnopqrstuvwxyz"' \
  -H 'Content-Type: application/zip' \
  --file-binary="@filename.zip" \
  https://my.daata.xyz/docs/spokes-platform

  curl -i -X POST -H "Content-Type: multipart/form-data"
  -F "data=@test.mp3;userid=1234" http://mysuperserver/media/upload/

  curl -i -X POST -H "Content-Type: multipart/form-data"
  -F "data=@test.mp3" http://mysuperserver/media/1234/upload/

  use -F needn't set "Content-Type: multipart/form-data"
*/

// Versioned is
func Versioned(w http.ResponseWriter, r *http.Request) error {
	return nil

}

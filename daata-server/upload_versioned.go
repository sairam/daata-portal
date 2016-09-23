package main

import "net/http"

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

// UploadVersioned is
func UploadVersioned(w http.ResponseWriter, r *http.Request) error {
	return nil

}

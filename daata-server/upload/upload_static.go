package upload

import "net/http"

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

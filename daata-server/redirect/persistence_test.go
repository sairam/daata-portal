package redirect

import "testing"

//
// type urlShortner struct {
// 	shortURL string
// 	longURL  string
// 	override bool
// }
//

func TestExists(t *testing.T) {
	url := &urlShortner{
		"test",
		"https://www.daata.xyz",
		false,
	}
	if url.exists() {
		t.Errorf("Redirect already exists")
	}
	url.CreateOrUpdate()
	if !url.exists() {
		t.Errorf("Unable to create redirect")
	}

}

package redirect

import (
	"testing"

	"../utils"
)

func TestExists(t *testing.T) {
	shortURL := utils.RandomString(4)
	// create a new urlShortner
	url := &urlShortner{
		shortURL,
		"https://www.daata.xyz",
		false,
	}
	if url.exists() {
		t.Errorf("Redirect already exists")
	}

	// persist
	url.insert()
	if !url.exists() {
		t.Errorf("Unable to create redirect")
	}

}

func TestUpdateWithoutInsert(t *testing.T) {
	shortURL := utils.RandomString(4)
	url := &urlShortner{
		shortURL,
		"https://www.daata.xyz",
		false,
	}
	url.update()

	rurl := &urlShortner{url.shortURL, "", false}
	rurl.read()
	if rurl.longURL != "" {
		t.Errorf("update should not be done with out insert")
	}

}

func TestUpdateWithInsert(t *testing.T) {
	shortURL := utils.RandomString(4)
	url := &urlShortner{
		shortURL,
		"https://www.daata.xyz",
		false,
	}
	url.insert()

	url.longURL = "https://daata.xyz"
	url.update()

	rurl := &urlShortner{url.shortURL, "", false}
	rurl.read()
	if rurl.longURL != "https://daata.xyz" {
		t.Errorf("Update of new url failed is %s", rurl.longURL)
	}

}

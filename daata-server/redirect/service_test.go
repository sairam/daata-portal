package redirect

import "testing"
import "../utils"

// test non existing url
func TestCreateOrUpdate(t *testing.T) {
	shortURL := utils.RandomString(4)
	url := &urlShortner{
		shortURL: shortURL,
		longURL:  "https://www.daata.xyz",
	}
	if url.CreateOrUpdate() != nil {
		t.Errorf("Create with a regular url failed")
	}
}

func TestCreateOrUpdateWithoutShortURL(t *testing.T) {
	url := &urlShortner{
		longURL: "https://www.daata.xyz",
	}
	if url.CreateOrUpdate() != nil {
		t.Errorf("Create with a non shortURL field FAILED")
	}

	if len(url.shortURL) != AutoGenerateShortURLLength {
		t.Errorf("Auto generated short URL size is not %d", AutoGenerateShortURLLength)
	}

}

// tests existing url with override
func TestCreateOrUpdateWithOverride(t *testing.T) {
	shortURL := utils.RandomString(4)
	url := &urlShortner{
		shortURL: shortURL,
		longURL:  "https://www.daata.xyz",
		override: true,
	}
	if url.CreateOrUpdate() != nil {
		t.Errorf("Creating a regular redirect FAILED")
	}

	url = &urlShortner{
		shortURL: shortURL,
		longURL:  "https://daata.xyz",
		override: true,
	}
	if url.CreateOrUpdate() != nil {
		t.Errorf("Create with an override url FAILED")
	}
}

// tests existing url without override. should fail
func TestCreateOrUpdateWithoutOverride(t *testing.T) {
	shortURL := utils.RandomString(4)
	url := &urlShortner{
		shortURL: shortURL,
		longURL:  "https://www.daata.xyz",
		override: false,
	}
	if url.CreateOrUpdate() != nil {
		t.Errorf("Creating a regular redirect FAILED")
	}

	url = &urlShortner{
		shortURL: shortURL,
		longURL:  "https://daata.xyz",
		override: false,
	}
	if url.CreateOrUpdate() == nil {
		t.Errorf("Create without an override url should have FAILED")
	}
}

package redirect

import (
	"errors"
	"regexp"
	"strings"

	"../utils"
)

// This interface for using the service
type serviceInterface interface {
	Validate() []error
	CreateOrUpdate() error
	Find() error
}

/* Service Layer */
// shortURL without protocol, just a string of [a-zA-Z0-9-/] are allowed
// longURL can be anything including http, https, itunes urls or any other.
// for now limit to http, https for regular users.
// leading and trailing slashes will be removed.

// CreateOrUpdateURL is the main method to add a new redirect
func (u *urlShortner) CreateOrUpdate() error {
	if u.shortURL == "" {
		u.shortURL = utils.RandomString(AutoGenerateShortURLLength)
	}

	var persistence persistenceInterface = u

	// insert on false
	// update on true

	if persistence.exists() {
		if u.override {
			return persistence.update()
		}
		return errors.New("URL already exists. use 'override' flag to replace")
	}
	return persistence.insert()
}

// add caching if required in service layer
func (u *urlShortner) Find() error {
	var persistence persistenceInterface = u
	if u.shortURL != "" && persistence.exists() {
		return persistence.read()
	}
	return errors.New("no such url exists")
}

// Validate verifies shortURL and longURL
func (u *urlShortner) Validate() []error {
	var errs []error
	var err []error

	err = validateURLs(u.shortURL, validateShortURL)
	if err != nil {
		errs = append(errs, err...)
	}

	err = validateURLs(u.longURL, validateBlankURL, validateLongURL)

	if err != nil {
		errs = append(errs, err...)
	}

	return errs
}

func validateURLs(url string, fs ...func(string) error) []error {
	var errs []error
	for _, f := range fs {
		err := f(url)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ensure there are no spaces, dots or any such
// whitelist with unicode chars.
// TODO - make a demo with emoticons
// validate if its a valid file system path
func validateShortURL(s string) error {
	if s == "" {
		return nil
	}
	pattern := `[\p{L}|\d_-]+`
	match, _ := regexp.MatchString(pattern, s)
	if match == true {
		return nil
	}
	return errors.New("unable to match pattern for short_url")
}

func validateBlankURL(s string) error {
	if s == "" {
		return errors.New("long_url is blank")
	}
	return nil
}
func validateLongURL(s string) error {
	if err := validateLongURLProto(s); err != nil {
		err1 := validateRelativePath(s)
		if err1 != nil {
			return err
		}
	}
	return nil
}

func validateRelativePath(s string) error {
	if s != "" && s[0] != '/' {
		return errors.New("long_url does not start with '/'")
	}
	return nil
}

func validateLongURLProto(s string) error {
	if !(strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")) {
		return errors.New("long_url is neither http or https protocol")
	}
	return nil
}

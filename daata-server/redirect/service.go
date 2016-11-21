package redirect

import (
	"errors"

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

func (u *urlShortner) Validate() []error {
	var errs []error
	var err []error

	err = validateURLs(u.shortURL, validateShortURL)
	if err != nil {
		errs = append(errs, err...)
	}

	err = validateURLs(u.longURL, validateBlankURL, validateLongURL, validateRelativePath)

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

func validateShortURL(_ string) error {
	// ensure there are no spaces, dots or any such
	// whitelist with unicode chars.
	// TODO - make a demo with emojicons
	// validate if its a valid file system path
	return nil
}

func validateBlankURL(str string) error {
	if str == "" {
		return errors.New("long_url is blank")
	}
	return nil
}

func validateRelativePath(str string) error {
	// if str[0] != '/' {
	// 	return errors.New("url does not start with '/'")
	// }
	// ensure does not have script tag
	return nil
}

func validateLongURL(_ string) error {
	// TODO check if url starts with http or https
	return nil
}

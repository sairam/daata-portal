package redirect

import "errors"

// TODO - check how to make this cleaner
func (u *urlShortner) Validate() error {
	var err error

	err = validateShortURL(u.shortURL)
	if err != nil {
		return err
	}

	err = validateBlankURL(u.longURL)
	if err != nil {
		return err
	}

	valid, err := validateLongURL(u.longURL)
	if err != nil {
		return err
	}

	if !valid {
		err = validateRelativePath(u.longURL)
		if err != nil {
			return err
		}
	}

	return nil
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
	if str[0] != '/' {
		return errors.New("url does not start with '/'")
	}
	// ensure does not have script tag
	return nil
}

func validateLongURL(_ string) (bool, error) {
	// TODO check if url starts with http or https
	return true, nil
}

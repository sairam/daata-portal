package redirect

import (
	"log"
	"os"
	"strings"
)

// Interface for persistence of data
type persistenceInterface interface {
	exists() bool
	read()
	insert() error
	update() error
}

// NOTE: persistence layer does not take care look at 'override' field from the model

/* Query Data Store */

// returns true if the file exists
// returns false if the file does not exist
func (u *urlShortner) exists() bool {
	if _, err := appFs.Stat(u.shortURL); os.IsNotExist(err) {
		return false
	}
	return true
}

func (u *urlShortner) read() error {
	data, err := fsutil.ReadFile(u.shortURL)
	if err != nil {
		return err
	}

	dataStr := string(data)
	u.longURL = strings.Split(dataStr, "\n")[0]

	return nil
}

func (u *urlShortner) insert() error {

	file, err := appFs.OpenFile(u.shortURL, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer file.Close()
	count, err := file.WriteString(u.longURL + "\n")
	log.Printf("no. of bytes written is %d\n", count)
	if err != nil {
		log.Printf("errror is %s", err)
	}
	return err
}

func (u *urlShortner) update() error {

	file, err := appFs.OpenFile(u.shortURL, os.O_WRONLY, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Truncate(0)
	file.WriteString(u.longURL + "\n")
	return nil
}

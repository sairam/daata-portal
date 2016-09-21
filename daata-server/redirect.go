package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type urlShortner struct {
	shortURL string
	longURL  string
}

type model struct {
	urlShortner
}

type urlShortnerForm struct {
	urlShortner
	Override bool
}

func query(shortURL string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT long_url FROM tiny_urls WHERE short_url = ?", shortURL).Scan(&longURL)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		// case err != nil:
		// 	return (longURL,err)
		// default:
		// 	return (longURL, nil)
	}
	return longURL, err

}

func insert(shortURL, longURL string) error {
	stmt, err := db.Prepare("INSERT INTO tiny_urls VALUES( ?, ? )")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(shortURL, longURL)

	defer stmt.Close()
	return err
}

func update(shortURL, longURL string) error {
	stmt, err := db.Prepare("UPDATE tiny_urls SET long_url = ? WHERE short_url = ? LIMIT 1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(longURL, shortURL)

	defer stmt.Close()
	return err
}

func setupMigrations() {
	_ = `
  CREATE TABLE tiny_urls (
   id int(11) unsigned NOT NULL AUTO_INCREMENT,
   short_url varchar(100) NOT NULL DEFAULT '',
   long_url tinytext NOT NULL,
   PRIMARY KEY (id),
   UNIQUE KEY short_url (short_url)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
  `
	// TODO: Add company_id, user_id who added this.
	// company_id is used for scope
}

// Setup is a one time function call to perform migrations and index updates
func init() {
	db, err := sql.Open("mysql", "root@/daata")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	// defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("Initialized DB connection")

	// Migrate the schema
	// db.AutoMigrate(&model{})

	setupMigrations()
}

// shortURL without protocol, just a string of [a-zA-Z0-9-/] are allowed
// longURL can be anything including http, https, itunes urls or any other.
// for now limit to http, https for regular users.
// leading and trailing slashes will be removed.

// insert query, don't update
func insertshortURL(shortURL, longURL string) error {
	makeEntryEvenIfExists(shortURL, longURL, false)
	return nil
}

// upsert query
func upsertshortURL(shortURL, longURL string) error {
	makeEntryEvenIfExists(shortURL, longURL, true)
	return nil
}

func makeEntryEvenIfExists(shortURL, longURL string, override bool) error {
	function := insert
	if existingURL, _ := query(shortURL); existingURL != "" {
		function = update
	}
	return function(shortURL, longURL)
}

func createNewURL(shortURL, longURL string, update bool) (string, error) {
	if shortURL == "" {
		shortURL = randomString(6)
	}
	if update {
		upsertshortURL(shortURL, longURL)
		// upsert query
	} else {
		insertshortURL(shortURL, longURL)
	}
	// TODO save into DB
	// update
	return "", nil
}

func findNewURL(str string) string {
	// TODO: query in DB
	return map[string]string{
		"test":   "https://www.google.com",
		"hello":  "http://www.hellobar.com",
		"yellow": "/yellow",
	}[str[1:]]
}

// Redirect is the main method which takes care of this functionality
// TODO Check Auth
func Redirect(w http.ResponseWriter, r *http.Request) {
	prefix := "/r"
	redirect := true
	if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
		fmt.Println(p)
		length := len(p) - 1
		if p[length] == '+' {
			redirect = false
			p = p[:length]
		}
		url := findNewURL(p)
		fmt.Println(url)
		fmt.Println(p)
		if url == "" {
			http.NotFound(w, r)
		} else {
			if redirect {
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			} else {
				fmt.Fprintf(w, url)
			}
		}
	} else {
		http.NotFound(w, r)
	}
}

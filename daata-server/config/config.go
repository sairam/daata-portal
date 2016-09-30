package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Config is a const
var Config = &AppConfig{}

// C returns
// TODO - Use Context to pass around config
// https://blog.golang.org/context
func C() *AppConfig {
	return Config
}

func init() {
	Config = loadConfig()
	fmt.Println("Loaded config")
}

func loadConfig() *AppConfig {
	f, err := os.Open("./config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	var config = &AppConfig{}

	if err := toml.Unmarshal(buf, config); err != nil {
		// panic(err)
	}
	return config
}

// Reload ..
// TODO - whenever the file changes we can reload the config
// Or can be done whenever we get a USR1 signal
func Reload() {
	Config = loadConfig()
}

// AppConfig ..
type AppConfig struct {
	Server      server
	Upload      upload
	Permissions uploadPermissions `toml:"upload_permissions"`
	Directories directories
	Redirect    redirect
}

type server struct {
	Port int
	URL  string
}

type upload struct {
	SizeLimit       int `toml:"size_limit"`
	Directory       string
	DirectoryLength int `toml:"directory_length"`
}

type uploadPermissions struct {
	Directory int
	File      int
}

type directories struct {
	Static  string
	Display string
}

type redirect struct {
	Length int
	Prefix string
}

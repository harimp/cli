package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// File is the default name of the JSON file where the config written.
	// The user can pass an alternate filename when using the CLI.
	// TODO: rename to .exercism.json
	File = ".exercism.go"
	// Host is the default hostname for fetching problems and submitting exercises.
	// TODO: We need to operate against two hosts (one for problems and one for submissions),
	// or define a proxy that both APIs can go through.
	Host = "http://exercism.io"
	// DemoDirname is the default directory to download problems to.
	DemoDirname = "exercism-demo"
)

// Config represents the settings for particular user.
// This defines both the auth for talking to the API, as well as
// where to put problems that get downloaded.
type Config struct {
	GithubUsername    string `json:"githubUsername"`
	APIKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
	Hostname          string `json:"hostname"`
}

// ToFile writes a Config to a JSON file.
func ToFile(path string, c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Your credentials have been written to %s\n", path)
	return nil
}

// FromFile loads a Config object from a JSON file.
func FromFile(path string) (c Config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &c)
	return
}

// HomeDir return's the user's canonical home directory.
// See: http://stackoverflow.com/questions/7922270/obtain-users-home-directory
// we can't cross compile using cgo and use user.Current()
func HomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
}

// Filename is the name of the JSON file containing the user's config.
func Filename(dir string) string {
	return fmt.Sprintf("%s/%s", dir, File)
}

// Demo is a default configuration for unauthenticated users.
func Demo() (c Config, err error) {
	demoDir, err := demoDirectory()
	if err != nil {
		return
	}
	c = Config{
		Hostname:          Host,
		APIKey:            "",
		ExercismDirectory: demoDir,
	}
	return
}

// ReplaceTilde replaces the short-hand home path with the absolute path.
func ReplaceTilde(oldPath string) string {
	return strings.Replace(oldPath, "~/", HomeDir()+"/", 1)
}

func demoDirectory() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return
	}
	dir = filepath.Join(dir, DemoDirname)
	return
}
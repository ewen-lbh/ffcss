package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/pelletier/go-toml"
)

// ReadFileBytes reads the content of ``filepath`` and returns the contents as a byte array
func ReadFileBytes(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b
}

// ReadFile reads the content of ``filepath`` and returns the contents as a string
func ReadFile(filepath string) string {
	return string(ReadFileBytes(filepath))
}

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

type ThemesList map[string]struct {
	repository string
	files      map[string][]string
	config     map[string]interface{}
}

func ReadThemesList() ThemesList {
	themesList := ThemesList{}
	doc := ReadFileBytes("themes.toml")
	toml.Unmarshal(doc, &themesList)
	return themesList
}

func Assert(t *testing.T, got interface{}, expected interface{}) {
	if got != expected {
		t.Errorf("expected: %s\ngot: %s", expected, got)
	}
}

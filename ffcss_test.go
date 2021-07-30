package main

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var currentUser user.User
var mockedHomedir string
var mockedProfile FirefoxProfile
var mockedStdout bytes.Buffer
var testarea = filepath.Join(cwd(), "testarea")

func init() {
	usr, _ := user.Current()
	currentUser = *usr
	mockedHomedir, _ = os.UserHomeDir()
	mockedProfile = NewFirefoxProfileFromPath(filepath.Join(mockedHomedir, ".mozilla", "firefox", "667ekipp.default-release"))
	out = &mockedStdout
}

// withUser Replaces %s with currentUser.Username in s
func withuser(s string) string {
	return fmt.Sprintf(s, currentUser.Username)
}

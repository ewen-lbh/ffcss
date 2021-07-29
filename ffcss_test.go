package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var currentUser user.User
var mockedHomedir string
var testarea = filepath.Join(cwd(), "testarea")

func init() {
	usr, _ := user.Current()
	currentUser = *usr
	mockedHomedir, _ = os.UserHomeDir()
}

// withUser Replaces %s with currentUser.Username in s
func withuser(s string) string {
	return fmt.Sprintf(s, currentUser.Username)
}

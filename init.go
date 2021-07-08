package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/MakeNowJust/heredoc"
	"github.com/docopt/docopt-go"
)

// RunCommandInit runs the command "init"
func RunCommandInit(args docopt.Opts) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	// TODO: set user{Chrome,Content,.js} by finding their path
	// TODO: only set assets if chrome/ actually exists
	content := fmt.Sprintf(heredoc.Doc(`
		ffcss: %d

		name: %s
		repository: %s

		userChrome: userChrome.css
		userContent: userContent.css
		user.js: user.js
		assets:
			- chrome/**
	`), VersionMajor, path.Dir(workingDir), getCurrentRepoRemote())

	err = ioutil.WriteFile(path.Join(workingDir, "ffcss.yaml"), []byte(content), 0700)

	return nil
}

// getCurrentRepoRemote returns the git repo's origin remote URL
// if any error occured while getting the URL, "# TODO" is returned
func getCurrentRepoRemote() string {
	var out bytes.Buffer
	command := exec.Command("git", "config", "--get", "remote.origin.url")
	command.Stdout = &out

	err := command.Run()
	if err != nil {
		fmt.Printf("WARNING: Could not get the current git remote origin's URL. Leaving repository entry blank.\n")
		return "# TODO"
	}
	return out.String()
}

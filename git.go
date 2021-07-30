package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// getCurrentRepoRemote returns the git repo's origin remote URL
// if any error occurred while getting the URL, the empty string is returned.
func getCurrentRepoRemote() string {
	var out bytes.Buffer
	command := exec.Command("git", "config", "--get", "remote.origin.url")
	command.Stdout = &out

	err := command.Run()
	if err != nil {
		warn("Could not get the current git remote origin's URL. Leaving repository entry blank.\n")
		return ""
	}
	return out.String()
}

func SwitchGitBranch(newBranch, clonedTo string) error {
	process := exec.Command("git", "switch", newBranch)
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

func SwitchGitCommit(commitSHA, clonedTo string) error {
	process := exec.Command("git", "checkout", commitSHA)
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

func SwitchGitTag(tagName, clonedTo string) error {
	process := exec.Command("git", "fetch", "--all", "--tags")
	process.Dir = clonedTo
	output, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}

	process = exec.Command("git", "checkout", "tags/"+tagName)
	process.Dir = clonedTo
	output, err = process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}

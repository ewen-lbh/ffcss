package main

import (
	"io/ioutil"
	"os"
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

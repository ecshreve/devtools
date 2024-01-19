package main

import "os"

// GetDiffFromFile reads a file containing a git diff and returns it as a string.
func GetDiffFromFile(filename string) (string, error) {
	fileContent, _ := os.ReadFile(filename)

	return string(fileContent), nil
}

package gencom

import (
	"os"
	"strings"
)

// GetDiffFromFile reads a file containing a git diff and returns it as a string.
func GetDiffFromFile(filename string) (string, error) {
	fileContent, _ := os.ReadFile(filename)

	return string(fileContent), nil
}

// foldString wraps the input string at the given column (72 in this case)
func foldString(s string, lineWidth int) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return wrapped
}

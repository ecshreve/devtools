package gencom

import (
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

// GetDiffFromFile reads a file containing a git diff and returns it as a string.
func GetDiffFromFile(filename string) (string, error) {
	fileContent, _ := os.ReadFile(filename)

	return string(fileContent), nil
}

// foldString wraps the input string at the given column (72 in this case)
func foldString(s string, lineWidth int) string {
	log.Info("foldString")
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

// checkRequiredCommands checks if the required commands are installed.
func checkRequiredCommands(cmds []string) {
	// Check if cmds are installed
	for _, cmd := range cmds {
		_, err := exec.LookPath(cmd)
		if err != nil {
			log.Fatal("command is not installed.", "cmd", cmd)
		}
	}
}

// checkOpenAIKey checks if the OPENAI_API_KEY environment variable is set.
func checkOpenAIKey() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set.")
	}
}

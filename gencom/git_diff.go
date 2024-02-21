package gencom

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

// GitClient is an interface that describes the methods
// that are used to interact with git.
type GitClient interface {
	GetDiff() (string, error)
}

// GitCommand is a struct that implements GitClient.
type GitCommand struct{}

// NewGitCommand creates a new instance of GitCommand
// and returns it as a GitClient. This is useful for
// mocking the GitCommand struct in other tests.
func NewGitCommand() GitClient {
	return GitCommand{}
}

// GetDiff gets the string representation of the git diff for
// the currently staged changes.
func (g GitCommand) GetDiff() (string, error) {
	log.Info("GitCommand.GetDiff")
	cmd := exec.Command("git", "diff", "--cached", "--unified=0")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	if len(output) == 0 {
		return "", fmt.Errorf("git diff returned empty")
	}

	log.Debug("GitCommand.GetDiff", "output", string(output), "err", err)
	return string(output), err
}

// MockGitCommand is a mock implementation of GitClient
// that can be used for testing.
//
// TODO: This is a bit of a hack. Revisit this.
type MockGitCommand struct {
	Diff string
}

// GetDiff gets the string representation of the git diff for
// for the mock.
func (mg MockGitCommand) GetDiff() (string, error) {
	log.Info("MockGitCommand.GetDiff")
	if mg.Diff == "" {
		return "", fmt.Errorf("mock git diff is empty")
	}

	return mg.Diff, nil
}

// ParseDiff takes a string containing a git diff and processes it.
// It returns a string containing the diff with the leading '+' removed.
func ParseDiff(diff string) string {
	log.Info("ParseDiff")
	var buf bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(diff))
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		// Filter lines: start with '+' but not '+++'
		if strings.HasPrefix(line, "+") {
			buf.WriteString(strings.TrimPrefix(line, "+"))
			buf.WriteString("\n")
			lineCount++
		}
		if strings.HasPrefix(line, "-") {
			buf.WriteString(strings.TrimPrefix(line, "+"))
			buf.WriteString("\n")
			lineCount++
		}
		if strings.HasPrefix(line, "@@") {
			buf.WriteString(strings.TrimPrefix(line, ""))
			buf.WriteString("\n")
			lineCount++
		}
		if lineCount >= 100 {
			break
		}
	}

	log.Debug("lines", "count", lineCount)
	return buf.String()
}

// summarizeDiff takes a string containing a git diff and processes it.
// It returns a string containing a cleaner version of the diff.
func summarizeDiff(diff string) string {
	log.Info("summarizeDiff")
	lines := strings.Split(diff, "\n")
	var summaryBuilder strings.Builder

	fileRegex := regexp.MustCompile(`^\+\+\+ b/(.*)`)
	changeRegex := regexp.MustCompile(`^@@ -(\d+)(,(\d+))? \+(\d+),(\d+) @@`)

	currentFile := ""
	changeLines := 0
	for _, line := range lines {
		if fileMatch := fileRegex.FindStringSubmatch(line); fileMatch != nil {
			currentFile = fileMatch[1]
			summaryBuilder.WriteString(fmt.Sprintf("\nFile Changed: %s\n", currentFile))
			changeLines = 0
		} else if changeMatch := changeRegex.FindStringSubmatch(line); changeMatch != nil {
			summaryBuilder.WriteString(fmt.Sprintf("Lines Removed: -%s, Lines Added: +%s\n", changeMatch[3], changeMatch[5]))
			changeLines = 0
		} else if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "---") {
			if len(strings.TrimSpace(strings.TrimPrefix(line, "+"))) == 0 {
				continue
			}

			if changeLines >= 10 {
				continue
			}
			summaryBuilder.WriteString(line + "\n")
			changeLines++

		}
	}

	log.Debug("summary", "summary", summaryBuilder.String())
	return summaryBuilder.String()
}

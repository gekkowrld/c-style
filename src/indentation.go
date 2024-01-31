package src

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func indentation() (bool, error) {
	var err error
	var err_msg displayStr
	str := string(fileInfo.FileContents)

	lines := strings.Split(str, "\n")

	for lineNumber, line := range lines {
		lineContent := line
		if hasLeadingSpaces(lineContent) {
			err_msg.Main = fmt.Sprintf("There is a leading space at %d", lineNumber)
			err_msg.Extra = fmt.Sprintf("Content:\n\t==> %s", lineContent)
			infoDisplay(err_msg)
		}
	}

	return true, err
}

func hasLeadingSpaces(lineContent string) bool {
	skipRegex := regexp.MustCompile(`[\t\v\f\n\r\x85\xA0]`)

	for _, r := range lineContent {
		if skipRegex.MatchString(string(r)) {
			continue
		} else if r == ' ' {
			return true
		} else if !unicode.IsSpace(r) {
			return false
		}
	}

	return false
}

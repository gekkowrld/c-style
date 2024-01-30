package src

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"unicode"
)

func indentation(filename string) (bool, error) {
	var err error
	var err_msg displayStr
	fileOPen, err := os.Open(filename)

	fileScan := bufio.NewScanner(fileOPen)
	fileScan.Split(bufio.ScanLines)

	lineNumber := 0

	for fileScan.Scan() {
		lineNumber++
		lineContent := fileScan.Text()
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

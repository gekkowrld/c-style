package src

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func indentation(filename string) (bool, error) {
	var err error
	var err_msg displayStr

	openFile, err := os.Open(filename)
	scanFile := bufio.NewScanner(openFile)
	scanFile.Split(bufio.ScanLines)

	line_number := 0
	indentLevel := 0
	inMultiLineComment := false
	inSwitchStatement := false

	for scanFile.Scan() {
		line_number++
		has_spaces := false
		lineContent := scanFile.Text()
		CommentContent := strings.TrimSpace(lineContent)

		if strings.HasPrefix(CommentContent, "/*") {
			inMultiLineComment = true
		}

		if inMultiLineComment && strings.HasSuffix(CommentContent, "*/") {
			inMultiLineComment = false
		}

		// Skip comments and whitespace (some).
		if inMultiLineComment {
			continue
		}
		if strings.HasPrefix(CommentContent, "//") {
			continue
		}
		skipRegex := regexp.MustCompile(`^[ \t\v\f\n\r]*$`)
		if skipRegex.MatchString(lineContent) {
			continue
		}

		if !inMultiLineComment {
			has_spaces = hasLeadingSpaces(lineContent)
			if has_spaces {
				err_msg.Main = fmt.Sprintf("Line %d contains a leading space, consider removing it or using tabs", line_number)
				errorDisplay(err_msg)
			}
			var foundIndent, expectedIndent int
			// Switch statements are the only ones with an exception.
			// This exception is only for case and all the indentation rules
			//  are followed accordingly
			// Use the CommentContent as switch is expected to be on its own line.
			switchRegex := regexp.MustCompile(`^[ \t]*switch`)
			if switchRegex.MatchString(lineContent) {
				// Check if this line contains the opening braces '{'
				// If not, pursue the next line(s) if they are whitespace, else
				//  don't bother, should have raised a syntax error in the compiler already.
				if strings.Contains(lineContent, "{") {
					continue
				}

				// Know that it is in a switch
				inSwitchStatement = true
			}

			if inSwitchStatement {
				expectedIndent, inSwitchStatement = handleSwitchIndent(lineContent, indentLevel)
				continue
			}

			foundIndent, expectedIndent = isIndentCorrect(lineContent, indentLevel)
			checkIndent := indentLevel
			if strings.Contains(lineContent, "}") {
				checkIndent--
			}
			if foundIndent != checkIndent {
				err_msg.Main = fmt.Sprintf("Expected %d tabs found %d tabs at line %d", indentLevel, foundIndent, line_number)
				errorDisplay(err_msg)
			}

			indentLevel = expectedIndent
		}

		if indentLevel > 3 {
			err_msg.Main = fmt.Sprintf("[%s]: You have have more than 3 indent levels (%d) at line %d, consider refactoring your code!", filename, indentLevel, line_number)
		}
	}

	return true, err
}

func handleSwitchIndent(lineContent string, indentLevel int) (int, bool) {
	stillInSwitch := true
	// Handle switch statement indentation
	if strings.Contains(lineContent, "case") || strings.Contains(lineContent, "default") {
		// Increase the indent level for case and default statements
		indentLevel++
	} else if strings.Contains(lineContent, "}") {
		// Decrease the indent level when the closing brace is encountered
		if indentLevel > 0 {
			indentLevel--
		}
		stillInSwitch = false
	}

	return indentLevel, stillInSwitch
}

func isIndentCorrect(lineContent string, expectedIndent int) (int, int) {
	foundIndent := 0

	// Get the number of tabs are in the current line.
	for _, fi := range lineContent {
		if fi == '\t' {
			foundIndent++
		} else if !unicode.IsSpace(fi) {
			break
		}
	}

	// Now try to know the expected indent.
	// Seems like indentation is broken in the version
	// of Betty that I was testing on
	// Even in the [Betty] tests, no indentation is actually enforced.
	/* Add/remove an indent level */
	for _, t := range lineContent {
		if t == '\x7B' {
			expectedIndent++
		} else if t == '\x7D' {
			if expectedIndent > 0 {
				expectedIndent--
			}
		}
	}

	return foundIndent, expectedIndent
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

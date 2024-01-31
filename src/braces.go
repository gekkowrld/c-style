package src

import (
	"fmt"
	"regexp"
	"strings"
)

func bracesPlacement() {

	var err_msg displayStr
	str := string(fileInfo.FileContents)

	lines := strings.Split(str, "\n")

	for line_number, lineContent := range lines {
		line_number++
		bracesRegex := regexp.MustCompile(`[\{\}]`)
		containsBrace := bracesRegex.MatchString(lineContent)
		if containsBrace {
			if isLoneBracket(lineContent) {
				continue
			}
			_, err := handleDoWhile(lineContent)
			if err != nil {
				err_msg.Main = fmt.Sprintf("At line %d -> %v", line_number, err)
				errorDisplay(err_msg)
				continue
			}

		}
	}
}

func isLoneBracket(lineContent string) bool {
	// First strip the line of all the whitespaces (front and back)
	replaceRegex := regexp.MustCompile(`[ \t\n\f\x0A\r]+`)
	replaceString := ""
	replacedString := replaceRegex.ReplaceAllString(lineContent, replaceString)

	bracesRegex := regexp.MustCompile(`[\{\}]`)
	isMatch := bracesRegex.MatchString(lineContent)

	if len(replacedString) == 1 && isMatch {
		return true
	}

	return false
}

func handleDoWhile(lineContent string) (bool, error) {
	var errMsg error
	doWhileLineRegex := regexp.MustCompile(`\bdo\s*\{`)
	if !doWhileLineRegex.MatchString(lineContent) && strings.Contains(lineContent, "do") {
		errMsg = fmt.Errorf("A 'do' should have an opening brace on the same line")
		return true, errMsg
	}

	// Check if 'while' is on the same line as the closing brace
	// Will never be called, should check up on it though
	endWhileLoopReg := regexp.MustCompile(`\}\s*while\b`)
	if !endWhileLoopReg.MatchString(lineContent) && strings.Contains(lineContent, "while") {
		errMsg = fmt.Errorf("A 'do-while' loop should have 'while' on the same line as the closing brace")
		return false, errMsg
	}

	return true, nil
}

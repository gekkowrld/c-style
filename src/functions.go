package src

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func handleFunction(filename string) {
	var err_msg displayStr

	// Read the file content
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
		return
	}

	funcDeclarationMatching := regexp.MustCompile(`[ \t]*[\w\*]+[\s\*]+[\w\*]+[\s]*\([\w\s\,\*]+\)[\s]?\{`)

	// State variables for tracking nested curly braces
	braceCount := 0
	startIndex := -1
	lineNumber := 0

	// Iterate through the file content
	var fileContentArray []string
	var foundAt []int
	for i, char := range string(fileContent) {
		if char == '\n' {
			lineNumber++
		}
		if char == '{' {
			if braceCount == 0 {
				// Found the start of a function
				startIndex = i
			}
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 {
				// Found the end of a function
				functionBody := string(fileContent[startIndex : i+1])
				fileContentArray = append(fileContentArray, functionBody)
				foundAt = append(foundAt, lineNumber)
			}
		}
	}

	declarationStr := funcDeclarationMatching.FindAll(fileContent, -1)
	declarationLen := len(declarationStr)

	funcInAFile := len(fileContentArray)
	for i := 0; i < funcInAFile; i++ {
		linesInAFunc := strings.Count(fileContentArray[i], "\n") - 1
		fromLocation := foundAt[i] - linesInAFunc
		toLocation := foundAt[i]
		if linesInAFunc > 40 {
			err_msg.Main = fmt.Sprintf("[%s]: Function has %d lines, from line %d to line %d", filename, linesInAFunc, fromLocation, toLocation)
			if declarationLen == funcInAFile {
				removeEndStr := regexp.MustCompile(`[\n\v\f\r\t\x85\xA0\v\{]+`)
				replacementStr := ""
				err_msg.Extra = fmt.Sprintf("==> Function declaration:\n\t%s", removeEndStr.ReplaceAll(declarationStr[i], []byte(replacementStr)))
			}
			infoDisplay(err_msg)
		}
	}
}

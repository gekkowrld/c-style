package src

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func checkHeader(filename string) {
	var err_msg displayStr
	f, err := os.Stat(filename)
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
	}
	if !strings.HasPrefix(filename, ".") && strings.HasSuffix(filename, ".h") && !f.IsDir() {
		checkForDoubleInclusion(filename)
		checkFunctiomDefinition(filename)
	}

}

func checkForDoubleInclusion(filename string) {
	var err_msg displayStr
	openFile, err := os.Open(filename)
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
	}

	isProtected := false
	endIfReg := regexp.MustCompile(`(?m)^[ \t]*\#endif`)
	ifndefReg := regexp.MustCompile(`(?m)^[ \t]*\#ifndef[ \t]{1,}[\S]+`)
	defineReg := regexp.MustCompile(`(?m)^[ \t]*\#define[ \t]{1,}[\S]+`)
	emptyLineReg := regexp.MustCompile(`^[\s]*$`)
	lineNumber := 0
	isIfDefCalled := false
	endifFound := 1
	firstEndIf := false

	scanFile := bufio.NewScanner(openFile)
	scanFile.Split(bufio.ScanLines)

	for scanFile.Scan() {
		lineContent := scanFile.Text()
		lineNumber++

		// Make sure that protection is off, else just don't flip it.
		if !isProtected && ifndefReg.MatchString(lineContent) {
			isIfDefCalled = true
			endifFound--
		}

		// Now turn on protection if all the required conditions are met.
		if !isProtected && isIfDefCalled && defineReg.MatchString(lineContent) {
			isProtected = true
			isIfDefCalled = false
		}
		// Make the subsequent lines unprotected
		if isProtected && endIfReg.MatchString(lineContent) {
			isProtected = false
		}

		// Count the endif(s)
		if endIfReg.MatchString(lineContent) {
			if endifFound == 0 {
				firstEndIf = true
			} else {
				firstEndIf = false
			}
			endifFound++
		}

		// If the user has more than one #endif
		if !isProtected && endIfReg.MatchString(lineContent) {
			if endifFound > 0 && !firstEndIf {
				err_msg.Main = fmt.Sprintf("You have extra #endif at line %d", lineNumber)
				infoDisplay(err_msg)
			}
		}

		// Emit errors
		if !isProtected && !isIfDefCalled && !firstEndIf && !emptyLineReg.MatchString(lineContent) {
			err_msg.Main = fmt.Sprintf("Line %d is not protected from double inclusion", lineNumber)
			infoDisplay(err_msg)
		}
	}
}

func checkFunctiomDefinition(filename string) {
	var err_msg displayStr
	bareFunctionReg := regexp.MustCompile(`(?m)^[ \t]*[\S]+[ ]{1,}[\S]+[ ]*\([ \S]*\)`)
	correctDeclarationReg := regexp.MustCompile(`(?m)^[ \t]*[\S]+[ ]{1,}[\S]+[ ]*\([ \S]*\)[ ]*;`)

	openFile, err := os.Open(filename)
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
	}

	fileScan := bufio.NewScanner(openFile)
	fileScan.Split(bufio.ScanLines)
	lineNumber := 0

	for fileScan.Scan() {
		lineContent := fileScan.Text()
		lineNumber++

		if bareFunctionReg.MatchString(lineContent) && !correctDeclarationReg.MatchString(lineContent) {
			err_msg.Main = fmt.Sprintf("On line %d: Only function declarations are allowed", lineNumber)
			infoDisplay(err_msg)
		}
	}
}

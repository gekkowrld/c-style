package src

import (
	"fmt"
	"regexp"
	"strings"
)

func bracesPlacement() {

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

		}
	}
      handleDoWhile()
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

func handleDoWhile() {
  var err_msg displayStr
  fileContent := string(fileInfo.FileContents)
  linesCont := strings.Split(fileContent, "\n")
	doWhileLineRegex := regexp.MustCompile(`\bdo\s*\{`)
	endWhileLoopReg := regexp.MustCompile(`\}\s*while\b`)

  isDoCalled := false
  lineNumber := 0

  for _, lineCont := range(linesCont){
    lineNumber++
    if doWhileLineRegex.MatchString(lineCont){
      isDoCalled = true
    }
    if isDoCalled && !endWhileLoopReg.MatchString(lineCont) && strings.Contains(lineCont, "while"){
      err_msg.Main = fmt.Sprintf("Closing while loop should be on the same line as the braces, mv while %d -> %d", lineNumber, (lineNumber-1))
      err_msg.Extra = fmt.Sprintf("Content:\n\t==> %s", lineCont)
      errorDisplay(err_msg)
      isDoCalled = false
    }
  }

}

package src

import (
	"fmt"
	"strings"
)

func checkLineLenght(lineLenght int) (bool, error) {
	var all_errors error

	str := string(fileInfo.FileContents)
	lines := strings.Split(str, "\n")

	for line_number, lineContent := range lines {
		line_lenght := len(lineContent)
		if line_lenght > lineLenght {
			var err_msg displayStr
			err_msg.Main = fmt.Sprintf("Line %d is longer than %d characters, it is %d characters", line_number, lineLenght, line_lenght)
			err_msg.Extra = fmt.Sprintf("Content:\n\t %s", lineContent)
			infoDisplay(err_msg)
		}
	}

	return true, all_errors
}

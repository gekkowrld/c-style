package src

import (
	"bufio"
	"fmt"
	"os"
)

func checkLineLenght(filename string, lineLenght int) (bool, error) {
	var all_errors error

	openFile, err := os.Open(filename)
	if err != nil {
		all_errors = err
	}

	scanFile := bufio.NewScanner(openFile)
	scanFile.Split(bufio.ScanLines)

	line_number := 0 // Keep track of the line number so as to display it

	for scanFile.Scan() {
		line_number++
		line_lenght := len(scanFile.Text())
		if line_lenght > lineLenght {
			var err_msg displayStr
			err_msg.Main = fmt.Sprintf("Line %d is longer than %d characters, it is %d characters", line_number, lineLenght, line_lenght)
			err_msg.Extra = fmt.Sprintf("Content:\n\t %s", scanFile.Text())
			infoDisplay(err_msg)
		}
	}

	return true, all_errors
}

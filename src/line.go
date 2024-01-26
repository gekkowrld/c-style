package src

import (
	"bufio"
	"fmt"
	"os"
)

func CheckLineLenght(filename string, lineLenght int) (bool, error)  {
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
      var err_msg Msg
      err_msg.Main = fmt.Sprintf("[%s:%d]:: Line longer than %d characters, it is %d characters", filename, line_number,lineLenght, line_lenght)
      err_msg.Extra = fmt.Sprintf("%s", scanFile.Text())
      InfoDisplay(err_msg)
    }
  }

  return true, all_errors
}


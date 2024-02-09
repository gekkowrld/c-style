package src

import (
	"fmt"
  "github.com/charmbracelet/lipgloss"
	"github.com/jwalton/go-supportscolor"
)

type displayStr struct {
	Main  string
	Extra string
}


func canSupportColour() bool {
	// Return false if the user has turned it off (and default)
	return !flagsPassed.Colour && supportscolor.Stdout().SupportsColor && supportscolor.Stdout().Has256
}

func displayMessage(prefix string, msg displayStr, colourType string) {
  var disp = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colourType))
	if flagsPassed.Verbose && msg.Extra != "" {
		if canSupportColour() {
      msg_disp := fmt.Sprintf("\n%s: %s\n%s\n", prefix, msg.Main, msg.Extra)
      fmt.Printf(disp.Render(msg_disp))
		} else {
			fmt.Printf("\n%s: %s\n%s", prefix, msg.Main, msg.Extra)
		}
	}
	if (!flagsPassed.Verbose || msg.Extra == "") && !flagsPassed.Quiet {
		if canSupportColour() {
      msg_disp := fmt.Sprintf("\n%s: %s",prefix, msg.Main)
			fmt.Printf(disp.Render(msg_disp))
		} else {
			fmt.Printf("\n%s: %s\n", prefix, msg.Main)
		}
	}
}

/* 
0 - grey/black (greyish)
1 - red
2 - blue
3 - yellow
4 - blue
5 - bright red (pinkish)
6 - greenish
7 - white
8 - grey/black (greyish)
9 - red
10 - greenish
*/

func infoDisplay(msg displayStr) {
	displayMessage("INFO", msg, "3")
  internalFlags.OutputCalled = true
}

func errorDisplay(msg displayStr) {
  displayMessage("ERROR", msg, "9")
  internalFlags.OutputCalled = true
}

func successDisplay(msg displayStr) {
  displayMessage("SUCCESS", msg, "6")
  internalFlags.OutputCalled = true
}

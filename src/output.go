
package src

import (
	"fmt"
	"github.com/fatih/color"
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

func displayMessage(prefix string, msg displayStr, colorFunc func(format string, a ...interface{})) {
	if flagsPassed.Verbose {
		if canSupportColour() {
			colorFunc("\n%s: %s\n%s\n", prefix, msg.Main, msg.Extra)
		} else {
			fmt.Printf("\n%s: %s\n%s\n", prefix, msg.Main, msg.Extra)
		}
	}
  if !flagsPassed.Verbose && !flagsPassed.Quiet {
		if canSupportColour() {
			colorFunc("\n%s: %s\n", prefix, msg.Main)
		} else {
			fmt.Printf("\n%s: %s\n", prefix, msg.Main)
		}
  }
}

func infoDisplay(msg displayStr) {
	displayMessage("INFO", msg, color.HiYellow)
}

func errorDisplay(msg displayStr) {
	displayMessage("ERROR", msg, color.HiRed)
}

func successDisplay(msg displayStr) {
	displayMessage("Success", msg, color.HiGreen)
}


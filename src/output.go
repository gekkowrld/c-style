package src

import (
	"fmt"
)

type displayStr struct {
  Main string
  Extra string
}

func infoDisplay(msg displayStr) { 
  if FlagsPassed.Verbose {
    fmt.Printf("\nINFO:\n%s\n%s\n", msg.Main, msg.Extra)
  }else if FlagsPassed.Quiet {
    // Display Nothing
  }else {
    fmt.Printf("INFO: %s\n", msg.Main)
  }
}

func errorDisplay(msg displayStr)  {
  if FlagsPassed.Verbose {
    fmt.Printf("\nERROR:\n%s\n%s\n", msg.Main, msg.Extra)
  }else if FlagsPassed.Quiet {
    // Display Nothing
  }else {
    fmt.Printf("ERROR: %s\n", msg.Main)
  }
}

func successDisplay(msg displayStr)  {
  if FlagsPassed.Verbose {
    fmt.Printf("\nSuccess:\n%s\n%s\n", msg.Main, msg.Extra)
  }else if FlagsPassed.Quiet {
    // Display Nothing
  }else {
    fmt.Printf("Success: %s\n", msg.Main)
  }}

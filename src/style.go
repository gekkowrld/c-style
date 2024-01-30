/*
Copyright © 2024 Gekko Wrld

*/

package src

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var flagsPassed struct {
	Verbose bool
	Quiet   bool
	LineLen int
	Colour  bool
	FuncLen int
}

var styleCmd = &cobra.Command{
	Use:   "style",
	Short: "Check if the code complies with the coding style",

	Run: func(cmd *cobra.Command, args []string) {
		var err_msg displayStr
		quietFlag := cmd.Flag("quiet").Changed
		verboseFlag := cmd.Flag("verbose").Changed
		recursiveFlag := cmd.Flag("recursive").Changed
		colourFlag := cmd.Flag("colour").Changed

		lineFlag := cmd.Flag("line")
		if lineFlag.Changed {
			lineLength, err := strconv.Atoi(lineFlag.Value.String())
			if err != nil {
				err_msg.Main = fmt.Sprintf("%v", err)
				errorDisplay(err_msg)
				return
			}
			flagsPassed.LineLen = lineLength
		} else {
			flagsPassed.LineLen = 80
		}
		funcFlag := cmd.Flag("func")
		if funcFlag.Changed {
			funcLength, err := strconv.Atoi(funcFlag.Value.String())
			if err != nil {
				err_msg.Main = fmt.Sprintf("%v", err)
				errorDisplay(err_msg)
				return
			}
			flagsPassed.FuncLen = funcLength
		} else {
			flagsPassed.FuncLen = 40
		}

		if quietFlag {
			flagsPassed.Quiet = true
		}
		if verboseFlag {
			flagsPassed.Verbose = true
		}
		if colourFlag {
			flagsPassed.Colour = true
		}

		argsPassed := len(args)
		if argsPassed > 0 {
			passedArg := args[0]
			passedArg, _ = filepath.Abs(passedArg)
			if unHiddenFileExists(passedArg) {
				callRelevantFunctions(passedArg)
			}

		} else if recursiveFlag {
			runRecursiveFlag()
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(styleCmd)
	styleCmd.PersistentFlags().BoolP("quiet", "q", false, "Don't show ANY output")
	styleCmd.PersistentFlags().BoolP("verbose", "v", false, "Display the output in verbose mode")
	styleCmd.PersistentFlags().BoolP("recursive", "r", false, "Run on all the c (.c and .h) files in the current directory and its children")
	styleCmd.PersistentFlags().Int("line", 80, "Change the line lenght to be used")
	styleCmd.PersistentFlags().BoolP("colour", "c", false, "Turn off colour display")
	styleCmd.PersistentFlags().Int("func", 40, "Change the function lenght to be used")
}

func runRecursiveFlag() {
	var err_msg displayStr
	currentDir, err := os.Getwd()
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
	}

	if !directoryExists(currentDir) {
		err_msg.Main = "Can't find the current working directory"
		errorDisplay(err_msg)
		os.Exit(1)
	}

	err = processFilesRecursively(currentDir)
	if err != nil {
		err_msg.Main = fmt.Sprintf("%v", err)
		errorDisplay(err_msg)
	}
}

func processFilesRecursively(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		fileExt := filepath.Ext(fullPath)
		requireFileExt := regexp.MustCompile(`^\.[hcHC]$`)

		if unHiddenFileExists(fullPath) && requireFileExt.MatchString(fileExt) {
			cwd, _ := os.Getwd()
			relative_path, _ := filepath.Rel(cwd, fullPath)
			if !flagsPassed.Quiet {
				fmt.Printf("======== %s ========\n", relative_path)
			}
			callRelevantFunctions(fullPath)
		}

		if directoryExists(fullPath) {
			err := processFilesRecursively(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func callRelevantFunctions(filename string) {
	indentation(filename)
	bracesPlacement(filename)
	checkLineLenght(filename, flagsPassed.LineLen)
	handleFunction(filename)
	handleComment(filename)
	checkHeader(filename)
}

func unHiddenFileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) || strings.HasPrefix(fileName, ".") {
		return false
	}
	return !info.IsDir()
}

func directoryExists(dirName string) bool {
	info, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

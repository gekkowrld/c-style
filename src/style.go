/*
Copyright © 2024 Gekko Wrld

*/

package src

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"crypto/sha1"

	"github.com/spf13/cobra"
)

var flagsPassed struct {
	Verbose           bool
	Quiet             bool
	LineLen           int
	Colour            bool
	FuncLen           int
	DisregardFileSize bool
}

var fileInfo struct {
	FileContents     []byte
	FileSize         int64
	FileName         string
	ReadErrors       error
	ModificationTime time.Time
	FileHash         []byte
}

var internalFlags struct {
	OutputCalled bool
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
		fileSizeFlag := cmd.Flag("bigfile").Changed

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
		if fileSizeFlag {
			flagsPassed.DisregardFileSize = true
		}

		argsPassed := len(args)
		if argsPassed > 0 {
			passedArg := args[0]
			passedArg, _ = filepath.Abs(passedArg)
			if unHiddenFileExists(passedArg) {
				// Set up the file info before doing anything.
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
	styleCmd.PersistentFlags().BoolP("bigfile", "b", false, "Ignore file size on large files")
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
			setFileInfo(fullPath)
			if !flagsPassed.Quiet {
				fmt.Printf("\n========[%x] Begin %s ========\n", fileInfo.FileHash[:5], relative_path)
			}
			callRelevantFunctions(fullPath)
			if !internalFlags.OutputCalled {
				var err_msg = displayStr{
					Main: "No problem found in here!, you are good to go!\n",
				}
				successDisplay(err_msg)
			}
			internalFlags.OutputCalled = false
			if !flagsPassed.Quiet {
				fmt.Printf("\n========[%x] End %s =======\n", fileInfo.FileHash[:4], relative_path)
			}
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
	indentation()
	bracesPlacement()
	checkLineLenght(flagsPassed.LineLen)
	handleFunction()
	handleComment()
	checkHeader()
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

func setFileInfo(filename string) error {
	var err_msg displayStr
	var err error
	if !unHiddenFileExists(filename) {
		err_msg.Main = fmt.Sprintf("%s doesn't exist, exiting with error 127", filename)
		errorDisplay(err_msg)
		return fmt.Errorf("file doesn't exist: %w", err)
	}
	fileContent, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fileContent.Close()

	_tmpfilename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	fileInfo.FileName = _tmpfilename

	stat, err := fileContent.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}
	fileInfo.FileSize = stat.Size()

	if fileInfo.FileSize > 1048576 && !flagsPassed.DisregardFileSize && fileInfo.FileSize < (2097152+2) {
		err_msg.Main = fmt.Sprintf("Your file is more than 1 MiB -> ( %d MiB )", (fileInfo.FileSize / (1024 * 1024)))
		infoDisplay(err_msg)
	}

	if fileInfo.FileSize < (2097152+2) || flagsPassed.DisregardFileSize {
		_tmp_content, err := os.ReadFile(filename)
		if err != nil {
			fileInfo.ReadErrors = err
		}

		fileInfo.FileContents = _tmp_content
	}

	if fileInfo.FileSize > (2097152+2) && !flagsPassed.DisregardFileSize {
		err_msg.Main = fmt.Sprintf("Your file is more than 2 MiB ( %d MiB ), exiting... \n\tpass --bigfile to continue anyway", (fileInfo.FileSize / (1024 * 1024)))
		errorDisplay(err_msg)
		os.Exit(2)
	}

	fileInfo.ModificationTime = stat.ModTime()

	newHash := sha1.New()
	if _, err := io.Copy(newHash, fileContent); err != nil {
		fileInfo.ReadErrors = err
	}
	fileInfo.FileHash = newHash.Sum(nil)

	return nil
}

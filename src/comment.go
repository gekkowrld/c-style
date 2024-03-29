package src

import (
	"fmt"
	"regexp"
	"strings"
)

func handleComment() {
	// First handle multiline comment (and single line).
	handleMultiLineComment()
}

func handleMultiLineComment() []string {
	var err_msg displayStr
	startComment := regexp.MustCompile(`(?m)^[ \t\v]*\/[\*]{1,}[\S ]+`)
	endComment := regexp.MustCompile(`(?m)^[ \t\w\*\-\d\S]*\*\/`)

	str := string(fileInfo.FileContents)

	lines := strings.Split(str, "\n")

	var commentLine []string
	var commentTemp map[int]string
	commentTemp = make(map[int]string)
	inAComment := false
	for lineNumber, lineContent := range lines {
		lineNumber++
		if startComment.Match([]byte(lineContent)) && !strings.Contains(lineContent, "*/") {
			inAComment = true
			commentLine = append(commentLine, lineContent)
			commentTemp[lineNumber] = lineContent
			start_c := strings.TrimSpace(lineContent)
			replaceStartStr := regexp.MustCompile(`(?mi)\/[\*]+`)
			start_c = replaceStartStr.ReplaceAllString(start_c, "")
			if len(strings.TrimSpace(start_c)) != 0 {
				err_msg.Main = fmt.Sprintf("Block comments use a leading /* on a separate line. mv %d -> %d?", lineNumber, (lineNumber + 1))
				err_msg.Extra = fmt.Sprintf("Content:\n\t==> %s", start_c)
				infoDisplay(err_msg)
			}
			continue
		}

		if inAComment {
			commentLine = append(commentLine, lineContent)
			commentTemp[lineNumber] = lineContent
		} else {
			handleSingleLineComment(lineContent, lineNumber)
		}

		if endComment.Match([]byte(lineContent)) && !strings.Contains(lineContent, "/*") {
			inAComment = false
			// Add a "novel" delimeter
			commentLine = append(commentLine, "-->====<--")

			end_c := strings.TrimSpace(lineContent)
			replaceStartStr := regexp.MustCompile(`(?mi)[\*]+\/`)
			end_c = replaceStartStr.ReplaceAllString(end_c, "")
			if len(strings.TrimSpace(end_c)) != 0 {
				err_msg.Main = fmt.Sprintf("Block comments use a trailing */ on a separate line. mv %d -> %d?", lineNumber, (lineNumber + 1))
				err_msg.Extra = fmt.Sprintf("Content:\n\t==> %s", end_c)
				infoDisplay(err_msg)
			}

			continue
		}

	}

	return commentLine
}

func handleSingleLineComment(lineContent string, lineNumber int) string {
	var comment_str string
	var err_msg displayStr
	var reCommentLine = regexp.MustCompile(`(?mU)\/[\*]{1,}[\S ]+\*\/`)
	var unCommentLine = regexp.MustCompile(`(?mi)[\/]{2,}[\S\s]*`)

	reGotStr := reCommentLine.Find([]byte(lineContent))
	if strings.Contains(string(reGotStr), "*/") {
		return string(reGotStr)
	}

	unGoStr := unCommentLine.Find([]byte(lineContent))
	if unGoStr != nil {
		err_msg.Main = fmt.Sprintf("C++ style comments are used at line %d", lineNumber)
		err_msg.Extra = fmt.Sprintf("Content:\n\t==> %s", unGoStr)
		infoDisplay(err_msg)
	}

	return comment_str
}

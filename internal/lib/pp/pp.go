package pp

import (
	"encoding/json"
	"regexp"
	"strings"
)

//-------------------------------------------------------------------------------------------------

func JSON(value any) string {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

//-------------------------------------------------------------------------------------------------

var (
	whitespaceOnly    = regexp.MustCompile("(?m)^[ \t]+$")
	leadingWhitespace = regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t\n])")
)

func Dedent(text string) string {
	var margin string

	if text == "" {
		return ""
	}

	if text[0] == '\n' {
		text = whitespaceOnly.ReplaceAllString(text[1:], "")
	} else {
		text = whitespaceOnly.ReplaceAllString(text, "")
	}

	indents := leadingWhitespace.FindAllStringSubmatch(text, -1)

	// Look for the longest leading string of spaces and tabs common to all
	// lines.
	for i, indent := range indents {
		if i == 0 {
			margin = indent[1]
		} else if strings.HasPrefix(indent[1], margin) {
			// Current line more deeply indented than previous winner:
			// no change (previous winner is still on top).
			continue
		} else if strings.HasPrefix(margin, indent[1]) {
			// Current line consistent with and no deeper than previous winner:
			// it's the new winner.
			margin = indent[1]
		} else {
			// Current line and previous winner have no common whitespace:
			// there is no margin.
			margin = ""
			break
		}
	}

	if margin != "" {
		text = regexp.MustCompile("(?m)^"+margin).ReplaceAllString(text, "")
	}
	return text
}

//-------------------------------------------------------------------------------------------------

package wordwrap

import (
	"strings"
	"unicode"
)

const esc = "\x1b["

func trimAnsiColor(s string) (ret string) {
	for ret = s; ; {
		startPos, endPos := rangeAnsiColor(ret)
		if startPos < 0 {
			return ret
		}
		if startPos == 0 {
			ret = ret[endPos:]
		} else {
			ret = ret[:startPos] + ret[endPos:]
		}
	}
}

func rangeAnsiColor(s string) (int, int) {
	escIndex := strings.Index(s, esc)
	if escIndex < 0 {
		return -1, -1
	}
	tmp := s[escIndex+len(esc):]
	// Invalid ansi format: "\x1b[<ascii character>"
	if len(tmp) <= 1 {
		return -1, -1
	}
	for i, c := range tmp {
		if c == 'm' {
			// returns start position and the length
			//
			// i.e. "hello\x1b[38;5;82m", length == 15
			// return 5, (5 + len(ansi) == 15)
			return escIndex, escIndex + (len(esc) + i + 1)
		}
		if !unicode.IsDigit(c) && c != ';' {
			break
		}
	}

	// Nothing 'm' period
	return -1, -1
}

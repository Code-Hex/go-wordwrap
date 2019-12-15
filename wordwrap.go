package wordwrap

import (
	"bytes"
	"unicode"

	runewidth "github.com/mattn/go-runewidth"
)

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint) string {
	limit := int(lim)
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current int
	var wordBuf, spaceBuf bytes.Buffer

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		if char == '\n' {
			if bufLen(wordBuf) == 0 {
				if current+bufLen(spaceBuf) > limit {
					current = 0
				} else {
					current += bufLen(spaceBuf)
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += bufLen(spaceBuf) + bufLen(wordBuf)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
			continue
		}

		if char == '\x1b' {
			ss := string(runes[i:])
			startPos, endPos := rangeAnsiColor(ss)
			if startPos >= 0 {
				tmp := ss[:endPos]
				wordBuf.WriteString(tmp)
				i += endPos - 1
				continue
			}
		}

		if unicode.IsSpace(char) {
			if bufLen(spaceBuf) == 0 || bufLen(wordBuf) > 0 {
				current += bufLen(spaceBuf) + bufLen(wordBuf)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

			spaceBuf.WriteRune(char)
			continue
		}

		l := runewidth.RuneWidth(char)
		if current+bufLen(wordBuf)+l > limit {
			wordBuf.WriteTo(buf)
			buf.WriteRune('\n')
			current = 0
			spaceBuf.Reset()
			wordBuf.Reset()
		} else if current+bufLen(spaceBuf)+bufLen(wordBuf)+l > limit && bufLen(wordBuf)+l < limit {
			buf.WriteRune('\n')
			current = 0
			spaceBuf.Reset()
		}
		wordBuf.WriteRune(char)
	}

	if wordBuf.Len() == 0 {
		if current+bufLen(spaceBuf) <= limit {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}

func bufLen(b bytes.Buffer) int {
	s := trimAnsiColor(b.String())
	return runewidth.StringWidth(s)
}

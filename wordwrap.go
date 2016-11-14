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
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer

	for _, char := range s {
		if char == '\n' {
			if buffLen(wordBuf) == 0 {
				if current+buffLen(spaceBuf) > lim {
					current = 0
				} else {
					current += buffLen(spaceBuf)
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += buffLen(spaceBuf) + buffLen(wordBuf)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if buffLen(spaceBuf) == 0 || buffLen(wordBuf) > 0 {
				current += buffLen(spaceBuf) + buffLen(wordBuf)
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

			spaceBuf.WriteRune(char)
		} else {

			wordBuf.WriteRune(char)

			if current+buffLen(spaceBuf)+buffLen(wordBuf) > lim && buffLen(wordBuf) < lim {
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
			} else if current+buffLen(wordBuf) == lim {
				wordBuf.WriteTo(buf)
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
				wordBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+uint(spaceBuf.Len()) <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}

func buffLen(b bytes.Buffer) uint {
	return uint(runewidth.StringWidth(b.String()))
}

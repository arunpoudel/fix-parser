package fix_parser

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r         *bufio.Reader
	separator string
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if ch == eof {
		return EOF, ""
	}

	if isTagValueSeparator(ch) {
		return TAGVALUESEPARATOR, string(ch)
	} else if s.isSeparator(ch) {
		return SEPARATOR, string(ch)
	}

	s.unread()
	return s.readTagOrValue()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) readTagOrValue() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	var token Token

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isTagValueSeparator(ch) {
			token = TAG
			s.unread()
			break
		} else if s.isSeparator(ch) {
			token = VALUE
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return token, buf.String()
}

func (s *Scanner) isSeparator(ch rune) bool {
	return ch == rune(s.separator[0])
}

func isTagValueSeparator(ch rune) bool {
	return ch == '='
}

func NewScanner(r io.Reader, separator string) *Scanner {
	return &Scanner{r: bufio.NewReader(r), separator: separator}
}

var eof = rune(0)

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

// Scan tokenizes a literal and return a Token and
// the literal that was tokenized. The token can be
// used to determine which part of the message we
// are at, eg: if the literal is "8" and the next rune
// after it is "=", then it tokenized the literal "8" is
// a tag. And if the literal is "=", then anything between "="
// and SEPARATOR is a value associated tag
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

	// Rollback the previously read rune to the buffer
	// as it is neither a tagvalue separator or a separator
	// which means that it is either a tag or a value, and
	// we don't want to miss the value
	s.unread()
	return s.readTagOrValue()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or eof is returned).
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
			// If the next character is "=" then this is
			// definitely a tag
			token = TAG
			s.unread()
			break
		} else if s.isSeparator(ch) {
			token = VALUE
			s.unread()
			break
		} else {
			// else it is a value associated with a tag
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

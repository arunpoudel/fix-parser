package fix_parser

import (
	"errors"
	"io"
	"strconv"
)

var (
	ErrMissingChecksum = errors.New("missing checksum")
	ErrEof             = errors.New("end of file")
)

type Parser struct {
	s                 *Scanner
	curTag            string
	bodyLength        int
	readBodyLength    int
	correctMessages   int
	incorrectMessages int
}

// Parse reads the message from io.Reader,
// parses it and returns a message or an error
func (p *Parser) Parse() (*Message, error) {
	m := NewMessage()
	p.bodyLength = 0
	p.readBodyLength = 0
	p.curTag = ""
	for {
		tok, lit := p.scan()
		if p.bodyLength != 0 {
			p.readBodyLength += len(lit)
		}
		if tok == EOF {
			return nil, ErrEof
		} else if tok == TAG {
			p.curTag = lit
		} else if tok == VALUE {
			if p.curTag == TAGBODYLENGTH {
				p.bodyLength, _ = strconv.Atoi(lit)
			}
			m.Add(p.curTag, lit)
		}
		m.AppendBuffer(lit)
		if p.bodyLength != 0 && p.readBodyLength > p.bodyLength {
			// Reached the end of message
			// Now just checksum is remaining
			// read it and assign it to the message
			tagTok, tagLit := p.scan() // Checksum TAG
			if tagTok != TAG || tagLit != TAGCHECKSUM {
				return nil, ErrMissingChecksum
			}
			//Skip separator
			_, _ = p.scan()
			// Get the checksum value
			valTok, valLit := p.scan()
			if valTok != VALUE {
				return nil, ErrMissingChecksum
			}
			m.Add(tagLit, valLit)
			_, _ = p.scan()
			break
		}
	}
	err := m.Validate()
	if err != nil {
		p.incorrectMessages += 1
		return nil, err
	}
	p.correctMessages += 1
	return m, err
}

// NumCorrectMessages returns the number of correct message read frm the reader
func (p *Parser) NumCorrectMessages() int {
	return p.correctMessages
}

// NumIncorrectMessages returns the number of incorrect message read frm the reader
func (p *Parser) NumIncorrectMessages() int {
	return p.incorrectMessages
}

func (p *Parser) scan() (tok Token, lit string) {
	return p.s.Scan()
}

// NewParser creates a new parser for FIX4.4 message parsing
func NewParser(r io.Reader, separator string) *Parser {
	return &Parser{s: NewScanner(r, separator)}
}

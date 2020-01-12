package fix_parser

// TODO: Write benchmarks

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

func (p *Parser) NumCorrectMessages() int {
	return p.correctMessages
}

func (p *Parser) NumIncorrectMessages() int {
	return p.incorrectMessages
}

func (p *Parser) scan() (tok Token, lit string) {
	return p.s.Scan()
}

func NewParser(r io.Reader, separator string) *Parser {
	return &Parser{s: NewScanner(r, separator)}
}

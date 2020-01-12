package fix_parser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrClippedTokenMissingValue = errors.New("clipped token missing value")
	ErrClippedTokenMissing      = errors.New("clipped token missing")
	ErrChecksumMismatch         = errors.New("checksum mismatch")
	tagToField                  = map[string]string{
		TAGVERSION:     "Version",
		TAGBODYLENGTH:  "BodyLength",
		TAGMSGTYPE:     "MsgType",
		TAGSENDINGTIME: "SendingTime",
		TAGSYMBOL:      "Symbol",
		TAGNOMDENTRIES: "NoMDEntries",
		TAGCHECKSUM:    "Checksum",
	}
)

type Header struct {
	Version     string
	BodyLength  int64
	MsgType     string
	SendingTime string
}

type Component struct {
	Symbol string
}

type Entry struct {
	EntryType string
	Price     float64
	Amount    float64
}

type Body struct {
	NoMDEntries int64
	curMdEntry  int64
	Entries     []Entry
}

type Tail struct {
	Checksum string
}

type Message struct {
	Header
	Component
	Body
	Tail
	clippedTag string
	buffer     string
}

func (m *Message) Validate() error {
	total := 0
	for _, r := range []rune(m.buffer) {
		total += int(r)
	}

	checksum := total % 256
	if fmt.Sprintf("%03d", checksum) != m.Checksum {
		return ErrChecksumMismatch
	}
	return nil
}

func (m *Message) AddEntries(tag string, value string) {
	if tag == TAGENTRYTYPE {
		m.curMdEntry += 1
		m.Entries[m.curMdEntry].EntryType = value
	} else if tag == TAGENTRYAMOUNT {
		f, _ := strconv.ParseFloat(value, 64)
		m.Entries[m.curMdEntry].Amount = f
	} else if tag == TAGENTRYPRICE {
		f, _ := strconv.ParseFloat(value, 64)
		m.Entries[m.curMdEntry].Price = f
	}
}

func (m *Message) Add(tag string, value string) {
	fieldName, ok := tagToField[tag]
	if !ok {
		m.AddEntries(tag, value)
		return
	}
	if reflect.ValueOf(m).Elem().FieldByName(fieldName).Kind() == reflect.String {
		reflect.ValueOf(m).Elem().FieldByName(fieldName).SetString(value)
	} else if reflect.ValueOf(m).Elem().FieldByName(fieldName).Kind() == reflect.Int64 {
		i, _ := strconv.ParseInt(value, 10, 64)
		reflect.ValueOf(m).Elem().FieldByName(fieldName).SetInt(i)
	} else if reflect.ValueOf(m).Elem().FieldByName(fieldName).Kind() == reflect.Float64 {
		f, _ := strconv.ParseFloat(value, 64)
		reflect.ValueOf(m).Elem().FieldByName(fieldName).SetFloat(f)
	}
	if tag == TAGNOMDENTRIES {
		i, _ := strconv.ParseInt(value, 10, 64)
		m.Entries = make([]Entry, i)
	}
}

func (m *Message) AppendBuffer(lit string) {
	m.buffer += lit
}

func NewMessage() *Message {
	return &Message{Body: Body{curMdEntry: -1}}
}

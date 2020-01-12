package fix_parser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	// The checksum provided in the message doesn't match the one generated
	ErrChecksumMismatch = errors.New("checksum mismatch")
	tagToField          = map[string]string{
		TAGVERSION:     "Version",
		TAGBODYLENGTH:  "BodyLength",
		TAGMSGTYPE:     "MsgType",
		TAGSENDINGTIME: "SendingTime",
		TAGSYMBOL:      "Symbol",
		TAGNOMDENTRIES: "NoMDEntries",
		TAGCHECKSUM:    "Checksum",
	}
)

// Header Part of FIX4.4 Message Format
type Header struct {
	Version     string
	BodyLength  int64
	MsgType     string
	SendingTime string
}

// Component Part of FIX4.4 Message Format
type Component struct {
	Symbol string
}

// Single MDEntry of FIX4.4 Message Format
type Entry struct {
	EntryType string
	Price     float64
	Amount    float64
}

// Body Part of FIX4.4 Message Format
type Body struct {
	NoMDEntries int64
	curMdEntry  int64
	Entries     []Entry
}

// Tail Part of FIX4.4 Message Format
type Tail struct {
	Checksum string
}

// FIX4.4 Partial message format
type Message struct {
	Header
	Component
	Body
	Tail
	clippedTag string
	buffer     string
}

// Validate validates if the checksum provided in the message
// matches the calculated checksum
// Our buffer contains message from the beginning until the
// start of checksum tag. The value of each char is added up
// and the reminder of the sum is divided by 256 (sum of all chars % 256)
// padded 3 chars is our checksum
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

// addEntries adds a MDEntry to the body
// As we only get one element at a give read pass
// We have to do it in multiple steps
func (m *Message) addEntries(tag string, value string) {
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

// Add adds a key => value pair to a FIX message
// Eg: If tag is TAGVERSION and value is FIX.4.4,
// then Message.Version is set to FIX.4.4
// We use reflect to set the value, which is
// supposed to be slower than having a switch case
// but when extending the model, this methods felt like
// easiest to achieve this cleanly without having a big
// switch case statement
func (m *Message) Add(tag string, value string) {
	fieldName, ok := tagToField[tag]
	if !ok {
		// If the tag is not in our list of tagToField map,
		// then it might be an entry, call the function and
		// if it is not, the function has no side effect
		m.addEntries(tag, value)
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
		// Init i number of entries as defined by TAGNOMDENTRIES
		i, _ := strconv.ParseInt(value, 10, 64)
		m.Entries = make([]Entry, i)
	}
}

// AppendBuffer appends a literal to the buffer
// which is later used to calculate the checksum of the message
func (m *Message) AppendBuffer(lit string) {
	m.buffer += lit
}

// NewMessage return new Message
func NewMessage() *Message {
	return &Message{Body: Body{curMdEntry: -1}}
}

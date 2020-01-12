package fix_parser

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	testParserTestTable := []struct {
		Name    string
		Message string
		Output  *Message
		Error   error
	}{
		{
			Name: "Test Valid Message",
			Message: "8=FIX.4.49=14235=W34=049=justtech52=20180206-21:43:36.00056=user262=TEST55=EURUSD268=2269=0270=1.31678271=100000.0269=1270=1.31667271=100000.010=057",
			Output: &Message{
				Header: Header{
					Version:     "FIX.4.4",
					BodyLength:  142,
					MsgType:     "W",
					SendingTime: "20180206-21:43:36.000",
				},
				Component: Component{
					Symbol: "EURUSD",
				},
				Body: Body{
					NoMDEntries: 2,
					Entries: []Entry{
						{
							EntryType: "0",
							Price:     1.31678,
							Amount:    100000,
						},
						{
							EntryType: "1",
							Price:     1.31667,
							Amount:    100000,
						},
					},
				},
				Tail: Tail{
					Checksum: "057",
				},
			},
			Error: nil,
		},
		{
			Name: "Test Invalid Message (Missing End Separator)",
			Message: "8=FIX.4.49=14235=W34=049=justtech52=20180206-21:43:36.00056=user262=TEST55=EURUSD268=2269=0270=1.31678271=100000.0269=1270=1.31667271=100000.010=057",
			Output: nil,
			Error:  ErrMissingChecksum,
		},
		{
			Name: "Test Invalid Message (Checksum Mismatch)",
			Message: "8=FIX.4.49=14235=W34=049=justtech52=20180206-21:43:36.00056=user262=TEST55=EURUSD268=2269=0270=1.31678271=100000.0269=1270=1.31667271=100000.010=054",
			Output: nil,
			Error:  ErrChecksumMismatch,
		},
	}

	sep := ""
	for _, test := range testParserTestTable {
		t.Run(test.Name, func(t *testing.T) {
			r := strings.NewReader(test.Message)
			p := NewParser(r, sep)
			m, err := p.Parse()
			if err != test.Error {
				t.Errorf("Expecting error to be %s, got %s", test.Error, err)
			}
			if test.Output != nil {
				if !cmp.Equal(test.Output, m, cmpopts.IgnoreUnexported(Message{}), cmpopts.IgnoreUnexported(Body{})) {
					t.Errorf("Wrong Message returned. Diff:")
					t.Errorf(
						cmp.Diff(test.Output, m, cmpopts.IgnoreUnexported(Message{}), cmpopts.IgnoreUnexported(Body{})),
					)
				}
			} else if m != nil {
				t.Errorf("Expecting parsed message to be empty. Got: %+v", m)
			}
		})
	}
}

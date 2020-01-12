package fix_parser

import (
	"fmt"
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	// TODO: Validate that the value we got from the function is correct
	f, err := os.Open("data/fix.bin")
	if err != nil {
		t.Fatalf("cannot load file, reason: %s", err.Error())
	}
	parser := NewParser(f, "")
	for {
		m, err := parser.Parse()
		if err == ErrEof {
			break
		}
		fmt.Println(fmt.Sprintf("Symbol: %s", m.Symbol))
		for _, entry := range m.Entries {
			if entry.EntryType == ENTRYTYPEBUY {
				fmt.Println(fmt.Sprintf("\tBuy Price: %f", entry.Price))
				fmt.Println(fmt.Sprintf("\tBuy Amount: %f", entry.Amount))
			} else {
				fmt.Println(fmt.Sprintf("\tSell Price: %f", entry.Price))
				fmt.Println(fmt.Sprintf("\tSell Amount: %f", entry.Amount))
			}
		}
		fmt.Println(fmt.Sprintf("Time and Date of Message: %s", m.SendingTime))
		fmt.Println("")
	}
}

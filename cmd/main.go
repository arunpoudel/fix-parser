package main

import (
	"fmt"
	fix "github.com/arunpoudel/fix-parser"
	"os"
)

func main() {
	f, err := os.Open("data/fix.bin")
	if err != nil {
		fmt.Println(fmt.Sprintf("cannot load file, reason: %s", err.Error()))
		return
	}
	parser := fix.NewParser(f, "")
	for {
		m, err := parser.Parse()
		if err == fix.ErrEof {
			break
		}
		fmt.Println(fmt.Sprintf("Symbol: %s", m.Symbol))
		for _, entry := range m.Entries {
			if entry.EntryType == fix.ENTRYTYPEBUY {
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

# fix-parser
Sample Parser for FIx 4.4

## How to Use
```
parser := fix_parser.NewParser(reader, "")
for {
		m, err := parser.Parse()
		if err == fix_parser.ErrEof {
			break
		} else if err != nil {
			// Skip any other error
            // Maybe log it?
			continue
		}
		fmt.Println(fmt.Sprintf("Symbol: %s", m.Symbol))
		for _, entry := range m.Entries {
			if entry.EntryType == fix_parser.ENTRYTYPEBUY {
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
```
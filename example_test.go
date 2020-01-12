package fix_parser_test

import (
	"fmt"
	fix_parser "github.com/arunpoudel/fix-parser"
	"os"
)

// This is an example explaining how to use the module
func Example_basic() {
	f, err := os.Open("data/fix.bin")
	if err != nil {
		fmt.Println(fmt.Sprintf("cannot load file, reason: %s", err.Error()))
		return
	}
	parser := fix_parser.NewParser(f, "")
	for {
		m, err := parser.Parse()
		if err == fix_parser.ErrEof {
			break
		} else if err != nil {
			// Skip any other error
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

	// Output:
	//Symbol: EURUSD
	//	Buy Price: 1.316780
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.316670
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.000
	//
	//Symbol: EURNOK
	//	Buy Price: 1.756350
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.756330
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.020
	//
	//Symbol: EURUSD
	//	Buy Price: 1.316770
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.316700
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.021
	//
	//Symbol: EURNOK
	//	Buy Price: 1.756370
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.756360
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.021
	//
	//Symbol: EURUSD
	//	Buy Price: 1.316830
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.316740
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.021
	//
	//Symbol: EURUSD
	//	Buy Price: 1.316740
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.316770
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.021
	//
	//Symbol: USDJPY
	//	Buy Price: 1.827830
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.827770
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.021
	//
	//Symbol: EURNOK
	//	Buy Price: 1.756390
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.756340
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.022
	//
	//Symbol: EURNOK
	//	Buy Price: 1.756340
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.756350
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.022
	//
	//Symbol: USDJPY
	//	Buy Price: 1.827900
	//	Buy Amount: 100000.000000
	//	Sell Price: 1.827820
	//	Sell Amount: 100000.000000
	//Time and Date of Message: 20180206-21:43:36.022
}

package fix_parser

type Token int

const (
	ILLEGAL Token = iota
	EOF
	EOM // End of message

	SEPARATOR
	TAG
	VALUE
	TAGVALUESEPARATOR

	TAGVERSION     string = "8"
	TAGBODYLENGTH  string = "9"
	TAGMSGTYPE     string = "35"
	TAGSENDINGTIME string = "52"
	TAGSYMBOL      string = "55"
	TAGNOMDENTRIES string = "268"
	TAGCHECKSUM    string = "10"
	TAGENTRYTYPE   string = "269"
	TAGENTRYPRICE  string = "270"
	TAGENTRYAMOUNT string = "271"
	ENTRYTYPEBUY   string = "0"
	ENTRYTYPESELL  string = "1"
)

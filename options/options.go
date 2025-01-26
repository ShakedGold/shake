package options

var Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Lexer   bool   `short:"l" long:"lexer" description:"Show the lexer output"`
	Parser  bool   `short:"p" long:"parser" description:"Show the parser output"`
	Input   string `short:"i" long:"input" description:"input shk file"`
}

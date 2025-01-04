package options

var Options struct {
	Verbose     []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Input       string `short:"i" long:"input" description:"input shk file"`
	Interperted bool   `short:"p" long:"interperted" description:"compile and run in interperted mode"`
}

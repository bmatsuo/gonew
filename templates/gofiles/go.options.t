{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import "flag" "fmt" "os" }}

// TODO Customize exported (capitalized) variables, types, and functions.

var (
    CmdHelpUsage string // Custom usage string.
    CmdHelpFoot  string // Printed after help.
)

// A struct that holds {{.Project.Target}}'s parsed command line flags.
type Options struct {
    Verbose bool
}

//  Create a flag.FlagSet to parse the {{.Project.Target}}'s flags.
func SetupFlags(opt *Options) *flag.FlagSet {
    fs := flag.NewFlagSet("{{.Project.Target}}", flag.ExitOnError)
    fs.BoolVar(&opt.Verbose, "v", false, "Verbose program output.")
    return setupUsage(fs)
}

// Check the {{.Project.Target}}'s flags and arguments for acceptable values.
// When an error is encountered, panic, exit with a non-zero status, or override
// the error.
func VerifyFlags(opt *Options, fs *flag.FlagSet) {
}

/**************************/
/* Do not edit below here */
/**************************/

//  Print a help message to standard error. See CmdHelpUsage and CmdHelpFoot.
func PrintHelp() { SetupFlags(&Options{}).Usage() }

//  Hook up CmdHelpUsage and CmdHelpFoot with flag defaults to function flag.Usage.
func setupUsage(fs *flag.FlagSet) *flag.FlagSet {
    printNonEmpty := func (s string) {
        if s != "" {
            fmt.Fprintf(os.Stderr, "%s\n", s)
        }
    }
    fs.Usage = func() {
        printNonEmpty(CmdHelpUsage)
        fs.PrintDefaults()
        printNonEmpty(CmdHelpFoot)
    }
    return fs
}

//  Parse the flags, validate them, and post-process (e.g. Initialize more complex structs).
func parseFlags() Options {
    var opt Options
    fs := SetupFlags(&opt)
    fs.Parse(os.Args[1:])
    VerifyFlags(&opt, fs)
    // Process the verified Options...
    return opt
}

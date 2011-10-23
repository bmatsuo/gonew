{{ template "go._package.t" $ }}
{{ template "go._head.t" $ }}
{{ import "flag" "fmt" "os" }}
/*
 *  Constants, variables, and functions that users may actually want to call
 *  are capitalized.
 */


var (
    // Set this variable to customize the help message header.
    // For example, `{{.Project.Target}} [options] action [arg2 ...]`.
    CommandLineHelpUsage string
    // Set this variable to print a message after the option specifications.
    // For example, "For more help:\n\t{{.Project.Target}} help [action]"
    CommandLineHelpFooter string
)

//  A struct that holds parsed option values.
//  TODO: Customize this struct with options for {{.Project.Target}}
type options struct {
    Verbose bool
}

//  Create a flag.FlagSet to parse the command line options/arguments.
//  TODO: Edit this function and add custom flags for {{.Project.Target}}
func setupFlags(opt *options) *flag.FlagSet {
    fs := flag.NewFlagSet("{{.Project.Target}}", flag.ExitOnError)
    fs.BoolVar(&(opt.Verbose), "v", false, "Verbose program output.")

    setupUsage(fs)
    return fs
}

//  Check the options for acceptable values. Panics or otherwise exits
//  with a non-zero exitcode when errors are encountered.
//  TODO: Make sure the {{.Project.Target}}'s flags are valid.
func verifyFlags(opt *options, fs *flag.FlagSet) {}

//  Print a help message to standard error. See constants CommandLineHelpUsage
//  and CommandLineHelpFooter.
func PrintHelp() {
    fs := setupFlags(&options{})
    fs.Usage()
}

//  Hook up the commandLineHelpUsage and commandLineHelpFooter strings
//  to the standard Go flag.Usage function.
func setupUsage(fs *flag.FlagSet) {
    printNonEmpty := func (s string) {
        if s != "" {
            fmt.Fprintf(os.Stderr, "%s\n", s)
        }
    }
    fs.Usage = func() {
        printNonEmpty(CommandLineHelpUsage)
        fs.PrintDefaults()
        printNonEmpty(CommandLineHelpFooter)
    }
}

//  Parse the command line options, validate them, and process them
//  further (e.g. Initialize more complex structs) if need be.
func parseFlags() options {
    var opt options
    fs := setupFlags(&opt)
    fs.Parse(os.Args[1:])
    verifyFlags(&opt, fs)
    // Process the verified options...
    return opt
}

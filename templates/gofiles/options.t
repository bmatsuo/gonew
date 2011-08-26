package main
/*
 *  Filename:    {{.file}}
 *  Author:      {{.name}} <{{.email}}>
 *  Created:     {{.date}}
 *  Description: Parse arguments and options from the command line.
 */
import (
    "os"
    "flag"
)

//  A struct that holds parsed option values.
type options struct {
    Verbose bool
}

//  Create a flag.FlagSet to parse the command line options/arguments.
func setupFlags(opt *options) *flag.FlagSet {
    var fs = flag.NewFlagSet("{{.gotarget}}", flag.ExitOnError)
    fs.BoolVar(&(opt.Verbose), "v", false, "Verbose program output.")
    return fs
}

//  Check the options for acceptable values. Panics or otherwise exits
//  with a non-zero exitcode when errors are encountered.
func verifyFlags(opt *options, fs *flag.FlagSet) {
}

//  Parse the command line options, validate them, and process them
//  further (e.g. Initialize more complex structs) if need be.
func parseFlags() options {
    var opt options
    var fs = setupFlags(&opt)
    fs.Parse(os.Args[1:])
    verifyFlags(&opt, fs)
    // Process the verified options...
    return opt
}

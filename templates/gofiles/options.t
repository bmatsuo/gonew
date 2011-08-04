package main
/*
 *  Filename:    {file}
 *  Author:      {name} <{email}>
 *  Created:     {date}
 *  Description: Parse arguments and options from the command line.
 */
import (
    "os"
    "flag"
)

//  A struct that holds parsed option values.
type options struct {.meta-left}
    verbose bool
{.meta-right}

//  Create a flag.FlagSet to parse the command line options/arguments.
func setupFlags(opt *options) *flag.FlagSet {.meta-left}
    var fs = flag.NewFlagSet("{gotarget}", flag.ExitOnError)
    fs.BoolVar(&(opt.verbose), "v", false, "Verbose program output.")
    return fs
{.meta-right}

//  Check the options for acceptable values. Panics or otherwise exits
//  with a non-zero exitcode when errors are encountered.
func verifyFlags(opt *options, fs *flag.FlagSet) {.meta-left}
{.meta-right}

//  Parse the command line options, validate them, and process them
//  further (e.g. Initialize more complex structs) if need be.
func parseFlags() options {.meta-left}
    var opt options
    var fs = setupFlags(&opt)
    fs.Parse(os.Args[1:])
    verifyFlags(&opt, fs)
    // Process the verified options...
    return opt
{.meta-right}

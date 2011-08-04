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

type options struct {.meta-left}
    verbose bool
{.meta-right}

func setupFlags(opt *options) *flag.FlagSet {.meta-left}
    var fs = flag.NewFlagSet("{gotarget}", flag.ExitOnError)
    fs.BoolVar(&(opt.verbose), "v", false, "Verbose program output.")
    return fs
{.meta-right}

func verifyFlags(opt *options, fs *flag.FlagSet) {.meta-left}
{.meta-right}

func ParseFlags() options {.meta-left}
    var opt options
    var fs = setupFlags(&opt)
    fs.Parse(os.Args[1:])
    verifyFlags(&opt, fs)
    // Process the verified options...
    return opt
{.meta-right}

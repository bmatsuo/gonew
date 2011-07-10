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

type Options struct {.meta-left}
    verbose bool
{.meta-right}

var opt = Options{.meta-left} {.meta-right}

func SetupFlags() *flag.FlagSet {.meta-left}
    var fs = flag.NewFlagSet("{gotarget}", flag.ExitOnError)
    fs.BoolVar(&(opt.verbose), "v", false, "Verbose program output.")
    return fs
{.meta-right}

func VerifyFlags(fs *flag.FlagSet) {.meta-left}
{.meta-right}

func ParseFlags() {.meta-left}
    var fs = SetupFlags()
    fs.Parse(os.Args[1:])
    VerifyFlags(fs)
    // Process the verified options...
{.meta-right}

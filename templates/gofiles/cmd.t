package main
/*
 *  Filename:    {{file}}
 *  Author:      {{name}} <{{email}}>
 *  Created:     {{date}}
 *  Description: {{desc}}
 *  Usage:       {{gotarget}} [options] ARGUMENT ...
 */
import (
    "os"
    "flag"
)

type Options struct {
    verbose bool
}
var opt = Options{}
func SetupFlags() *flag.FlagSet {
    var fs = flag.NewFlagSet("{{gotarget}}", flag.ExitOnError)
    fs.BoolVar(&(opt.verbose), "v", "Verbose program output.")
    return fs
}
func VerifyFlags() {
}
func ParseFlags() {
    var fs = SetupFlags()
    ps.Parse(os.Args[1:])
}

func main() {
    ParseFlags()
}

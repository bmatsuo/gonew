package main
/*
 *  Filename:    template.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 17:55:40 PDT 2011
 *  Description: 
 */
import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "path/filepath"
    "github.com/hoisie/mustache.go"
)

func WriteTemplate(filename, desc string, dict map[string]string, path...string) os.Error {
    var (
        tpath = filepath.Join(path...)
        template = mustache.RenderFile(tpath, dict, map[string]string{"file":filename})
    )
	if DEBUG || VERBOSE {
		fmt.Printf("Creating %s %s", desc, filename)
    }
    if DEBUG_LEVEL > 0 {
        log.Printf("template: %s", tpath)
        if DEBUG_LEVEL > 1 {
            log.Printf("context:\n%v", dict)
            if DEBUG_LEVEL > 2 {
                log.Print("\n", template, "\n")
            }
        }
    }
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(filename, templout, FilePermissions)
    return errWrite
    return nil
}

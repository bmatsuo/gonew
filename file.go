package main
/*
 *  Filename:    file.go
 *  Package:     main
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sun Jul  3 16:57:42 PDT 2011
 *  Description: 
 */
import (
    "os"
    "log"
    "io/ioutil"
    "path/filepath"
    "github.com/hoisie/mustache.go"
)

type File struct {
    Name string
    User string
    Pkg  string
    Repo RepoType
    Host RepoHost
}

func (f File) Create() os.Error {
    var (
        lib = f.Name + ".go"
        dict = map[string]string{
            "file":lib,
            "name":AppConfig.Name,
            "email":AppConfig.Email,
            "date":DateString(),
            "year":YearString(),
            "gotarget":f.Pkg}
        tpath = filepath.Join(GetTemplateRoot(), "gofiles", "lib.t")
        template = mustache.RenderFile(tpath, dict)
    )
	if DEBUG {
		log.Printf("Creating library %s", lib)
        if DEBUG_LEVEL > 0 {
		    log.Printf("    template: %s", tpath)
            if DEBUG_LEVEL > 1 {
		        log.Print("\n", template, "\n")
            }
        }
	}
    var templout = make([]byte, len(template))
    copy(templout, template)
    var errWrite = ioutil.WriteFile(lib, templout, FilePermissions)
    return errWrite
}

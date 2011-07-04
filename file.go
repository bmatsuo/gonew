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
)

type File struct {
    Name string
    User string
    Pkg  string
    Repo RepoType
    Host RepoHost
}

func (f File) GenerateDictionary() map[string]string {
    var (
        lib = f.Name + ".go"
        dict = map[string]string{
            "file":lib,
            "name":AppConfig.Name,
            "email":AppConfig.Email,
            "date":DateString(),
            "year":YearString(),
            "gotarget":f.Pkg}
    )
    return dict
}

func (f File) Create() os.Error {
    var (
        dict = f.GenerateDictionary()
        errWrite = WriteTemplate(dict["file"], "library", dict,
                "gofiles", "lib.t")
    )
    if errWrite != nil {
        return errWrite
    }
    // TODO: check the new file into git under certain conditions...
    return nil
}

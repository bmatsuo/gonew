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
    Name   string
    User   string
    Pkg    string
    Repo   RepoType
    Host   RepoHost
}
type TestFile struct {
    Name   string
    Pkg    string
    Repo   RepoType
    Host   RepoHost
}

func (t TestFile) GenerateDictionary() map[string]string {
    var (
        test = t.Name + "_test.go"
        dict = map[string]string{
            "file":test,
            "name":AppConfig.Name,
            "email":AppConfig.Email,
            "date":DateString(),
            "year":YearString(),
            "gotarget":t.Pkg}
    )
    return dict
}
func (t TestFile) TemplatePath() []string {
    return []string{"testfiles", "pkg.t"}
}
func (t TestFile) Create() os.Error {
    var (
        dict = t.GenerateDictionary()
        errWrite = WriteTemplate(dict["file"], "library", dict, t.TemplatePath()...)
    )
    if errWrite != nil {
        return errWrite
    }
    // TODO: check the new file into git under certain conditions...
    return nil
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
func (f File) TemplatePath() []string {
    return []string{"gofiles", "lib.t"}
}
func (f File) Create() os.Error {
    var (
        dict = f.GenerateDictionary()
        errWrite = WriteTemplate(dict["file"], "library", dict, f.TemplatePath()...)
    )
    if errWrite != nil {
        return errWrite
    }
    // TODO: check the new file into git under certain conditions...
    // Create a test for the new file.
    var (
        test = f.TestFile()
        errTestCreate = test.Create()
    )
    if errTestCreate != nil {
        return errTestCreate
    }
    return nil
}
func (f File) TestFile() TestFile {
    return TestFile{Name:f.Name, Pkg:f.Pkg, Repo:f.Repo, Host:f.Host}
}

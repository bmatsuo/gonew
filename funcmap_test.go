package main

/*  Filename:    funcmap_test.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Sun Oct 23 01:45:33 PDT 2011
 *  Description: For testing funcmap.go
 */

import (
    "testing"
    "text/template"
    "bytes"
)

var Fn = DefaultFuncMap()

func RunTemplateTest(name, tstr string, data interface{}, expectation string, T *testing.T) {
    t, err := template.New(name).Funcs(Fn).Parse(tstr)
    if err != nil {
        T.Fatalf("Template %s did not parse; %s", name, err.String())
    }
    b := new(bytes.Buffer)
    err = t.Execute(b, data)
    if err != nil {
        T.Fatalf("Template %s execution error; ", name, err.String())
    }
    produced := string(b.Bytes())
    if produced != expectation {
        T.Errorf("Executed template differs from expectation\n\t'%s'\n\t'%s'", produced, expectation)
    }
}

func TestImport(T *testing.T) {
    RunTemplateTest("TestImport-Nothing", "{{import}}", nil, "import ()", T)
    RunTemplateTest(
        "TestImport-Single",
        `{{import "abc"}}`, nil,
        "import (\n    \"abc\"\n)",
        T)
    RunTemplateTest(
        "TestImport-Multi",
        `{{import "abc" "def"}}`, nil,
        "import (\n    \"abc\"\n    \"def\"\n)",
        T)
    RunTemplateTest(
        "TestImport-complex",
        `{{import . "abc"}}`, []interface{}{"fmt", "os"},
        "import (\n    \"fmt\"\n    \"os\"\n    \"abc\"\n)",
        T)
}

func TestGoFunction(T *testing.T) {
    RunTemplateTest(
        "TestGoFunction-main",
        `{{func "main"}}`, nil,
        "func main() {\n\n}",
        T)
    RunTemplateTest(
        "TestGoFunction-test",
        `{{func "TestOne" "T *testing.T"}}`, nil,
        "func TestOne(T *testing.T) {\n\n}",
        T)
}

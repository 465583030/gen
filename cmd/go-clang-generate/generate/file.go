package generate

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"golang.org/x/tools/imports"
)

// File represents a generation file
type File struct {
	Name string

	HeaderFiles map[string]struct{}

	Functions []interface{}
	Enums     []*Enum
	Structs   []*Struct
}

// newFile creates a new blank file
func newFile(name string) *File {
	return &File{
		Name: name,

		HeaderFiles: map[string]struct{}{},
	}
}

var templateGenerateFile = template.Must(template.New("go-clang-generate-file").Parse(`package phoenix

{{range $h, $dunno := $.HeaderFiles}}
// #include "{{$h}}"
{{end}}// #include "go-clang.h"
import "C"

{{range $i, $f := $.Functions}}
{{$f}}
{{end}}

{{range $i, $e := $.Enums}}
{{$e.Comment}}
type {{$e.Name}} {{$e.UnderlyingType}}

const (
{{range $i, $ei := .Items}}	{{if $ei.Comment}}{{$ei.Comment}}
	{{end}}{{$ei.Name}}{{if eq $i 0}} {{$e.Name}}{{end}} = C.{{$ei.CName}}
{{end}}
)

{{range $i, $m := $e.Methods}}
{{$m}}
{{end}}
{{end}}

{{range $i, $s := $.Structs}}
{{$s.Comment}}
type {{$s.Name}} struct {
	c {{if $s.IsPointerComposition}}*{{end}}C.{{if not $s.CNameIsTypeDef}}struct_{{end}}{{$s.CName}}
}
{{range $i, $m := $s.Methods}}
{{$m}}
{{end}}
{{end}}
`))

func (f *File) generate() error {
	for _, e := range f.Enums {
		if e.HeaderFile != "" {
			f.HeaderFiles[e.HeaderFile] = struct{}{}
		}
	}
	for _, s := range f.Structs {
		if s.HeaderFile != "" {
			f.HeaderFiles[s.HeaderFile] = struct{}{}
		}
	}

	var b bytes.Buffer
	if err := templateGenerateFile.Execute(&b, f); err != nil {
		return err
	}

	filename := f.Name + "_gen.go"

	bo := b.Bytes()
	out, err := imports.Process(filename, bo, nil)
	if err != nil {
		// Write the file anyway so we can look at the problem
		if err := ioutil.WriteFile(filename, bo, 0600); err != nil {
			return err
		}

		return err
	}

	return ioutil.WriteFile(filename, out, 0600)
}
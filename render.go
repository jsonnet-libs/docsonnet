package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var indexTmpl = strings.Replace(`
# package {{.Name}}
´´´jsonnet
local {{.Name}} = import "{{.Import}}"
´´´

{{.Help}}

## Index

{{ range .Index }}{{ $l := mul .Level 2}}{{repeat (int $l) " "}}* ´{{ .Line }}´
{{ end }}

{{ range .Fields }} {{.Render}}
{{ end }}

`, "´", "`", -1)

type Renderable interface {
	Render() string
}

var objTmpl = strings.Replace(`
## {{ .Name }}
{{ .Help }}
`, "´", "`", -1)

type obj struct {
	Object
}

func (o obj) Render() string {
	t := template.Must(template.New("").Parse(objTmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, o); err != nil {
		panic(err)
	}
	return buf.String()
}

type fn struct {
	Function
}

func (f fn) Params() string {
	args := make([]string, 0, len(f.Args))
	for _, a := range f.Args {
		arg := a.Name
		if a.Default != nil {
			arg = fmt.Sprint("%s=%v", arg, a.Default)
		}
		args = append(args, arg)
	}

	return strings.Join(args, ", ")
}

func (f fn) Signature() string {
	return fmt.Sprintf("%s(%s)", f.Name, f.Params())
}

var fnTmpl = strings.Replace(`
### fn {{ .Name }}
´´´
{{ .Signature }}
´´´
{{ .Help }}
`, "´", "`", -1)

func (f fn) Render() string {
	t := template.Must(template.New("").Parse(fnTmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, f); err != nil {
		panic(err)
	}
	return buf.String()
}

func render(d Doc) (string, error) {
	tmpl := template.Must(template.New("").Funcs(sprig.TxtFuncMap()).Parse(indexTmpl))
	if err := tmpl.Execute(os.Stdout, &doc{
		Name:   d.Name,
		Import: d.Import,
		Help:   d.Help,
		Index:  buildIndex(d.API, 0),
		Fields: renderables(d.API, ""),
	}); err != nil {
		return "", err
	}

	return "", nil
}

type doc struct {
	Name   string
	Import string
	Help   string

	Index  []indexElem
	Fields []Renderable
}

type indexElem struct {
	Line  string
	Level int
}

func renderables(fields map[string]Field, prefix string) []Renderable {
	rs := []Renderable{}
	for _, f := range fields {
		switch {
		case f.Function != nil:
			fnc := fn{*f.Function}
			fnc.Name = strings.TrimPrefix(prefix+"."+fnc.Name, ".")
			rs = append(rs, fnc)
		case f.Object != nil:
			o := obj{*f.Object}
			o.Name = strings.TrimPrefix(prefix+"."+o.Name, ".")
			rs = append(rs, o)

			childs := renderables(o.Fields, o.Name)
			rs = append(rs, childs...)
		}
	}
	return rs
}

func buildIndex(fields map[string]Field, level int) []indexElem {
	elems := []indexElem{}
	for _, f := range fields {
		line := indexLine(f)
		elems = append(elems, indexElem{
			Line:  line,
			Level: level,
		})

		if f.Object != nil {
			childs := buildIndex(f.Object.Fields, level+1)
			elems = append(elems, childs...)
		}
	}
	return elems
}

func indexLine(f Field) string {
	switch {
	case f.Function != nil:
		return "fn " + fn{*f.Function}.Signature()
	case f.Object != nil:
		return fmt.Sprintf("obj %s", f.Object.Name)
	}
	panic("wtf")
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-clix/cli"
	"github.com/google/go-jsonnet"
	"github.com/markbates/pkger"
)

type Package struct {
	Name   string `json:"name"`
	Import string `json:"import"`
	Help   string `json:"help"`

	API Fields             `json:"api,omitempty"`
	Sub map[string]Package `json:"sub,omitempty"`
}

func main() {
	root := &cli.Command{
		Use:   "docsonnet",
		Short: "Utility to parse and transform Jsonnet code that uses the docsonnet extension",
	}

	root.AddCommand(loadCmd(), renderCmd())

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func loadCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "load",
		Short: "extracts docsonnet from Jsonnet and prints it as JSON",
		Args:  cli.ArgsExact(1),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		pkg, err := Load(args[0])
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return cmd
}

func renderCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "render",
		Short: "writes all found docsonnet packages to Markdown (.md) files, suitable for e.g. GitHub",
		Args:  cli.ArgsExact(1),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		pkg, err := Load(args[0])
		if err != nil {
			return err
		}

		fmt.Println(render(*pkg))
		return nil
	}

	return cmd
}

func Load(filename string) (*Package, error) {
	file, err := pkger.Open("/load.libsonnet")
	if err != nil {
		return nil, err
	}
	load, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	vm := jsonnet.MakeVM()
	importer, err := newImporter()
	if err != nil {
		return nil, err
	}
	vm.Importer(importer)

	vm.ExtCode("main", fmt.Sprintf(`(import "%s")`, filename))
	data, err := vm.EvaluateSnippet("load.libsonnet", string(load))
	if err != nil {
		return nil, err
	}

	var d Package
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatalln(err)
	}

	return &d, nil
}

type Importer struct {
	fi   jsonnet.FileImporter
	util jsonnet.Contents
}

func newImporter() (*Importer, error) {
	file, err := pkger.Open("/doc-util/main.libsonnet")
	if err != nil {
		return nil, err
	}
	load, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &Importer{
		fi:   jsonnet.FileImporter{},
		util: jsonnet.MakeContents(string(load)),
	}, nil
}

func (i *Importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	if importedPath == "doc-util/main.libsonnet" {
		return i.util, "<internal>", nil
	}

	return i.fi.Import(importedFrom, importedPath)
}

type Object struct {
	Name string `json:"-"`
	Help string `json:"help"`

	// children
	Fields Fields `json:"fields"`
}

type Fields map[string]Field

func (fPtr *Fields) UnmarshalJSON(data []byte) error {
	if *fPtr == nil {
		*fPtr = make(Fields)
	}
	f := *fPtr

	tmp := make(map[string]Field)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for k, v := range tmp {
		switch {
		case v.Function != nil:
			v.Function.Name = k
		case v.Object != nil:
			v.Object.Name = k
		case v.Value != nil:
			v.Value.Name = k
		}
		f[k] = v
	}

	return nil
}

// Field represents any field of an object.
type Field struct {
	// Function value
	Function *Function `json:"function,omitempty"`
	// Object value
	Object *Object `json:"object,omitempty"`
	// Any other value
	Value *Value `json:"value,omitempty"`
}

func (o *Field) UnmarshalJSON(data []byte) error {
	type fake Field

	var f fake
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	switch {
	case f.Function != nil:
		o.Function = f.Function
	case f.Object != nil:
		o.Object = f.Object
	case f.Value != nil:
		o.Value = f.Value
	default:
		return errors.New("field has no value")
	}

	return nil
}

func (o Field) MarshalJSON() ([]byte, error) {
	if o.Function == nil && o.Object == nil && o.Value == nil {
		return nil, errors.New("field has no value")
	}

	type fake Field
	return json.Marshal(fake(o))
}

type Function struct {
	Name string `json:"-"`
	Help string `json:"help"`

	Args []Argument `json:"args,omitempty"`
}

type Type string

const (
	TypeString = "string"
	TypeNumber = "number"
	TypeBool   = "boolean"
	TypeObject = "object"
	TypeArray  = "array"
	TypeAny    = "any"
	TypeFunc   = "function"
)

type Value struct {
	Name string `json:"-"`
	Help string `json:"help"`

	Type Type `json:"type"`
}

type Argument struct {
	Type    Type        `json:"type"`
	Name    string      `json:"name"`
	Default interface{} `json:"default"`
}

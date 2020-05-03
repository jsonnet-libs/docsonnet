package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/sh0rez/docsonnet/pkg/docsonnet"
)

func main() {
	data, err := eval()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))

	var d DS
	if err := json.Unmarshal(data, &d); err != nil {
		log.Fatalln(err)
	}

	pkg := load(d)
	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(out))
}

func load(d DS) docsonnet.Package {
	pkg := d.Package()
	pkg.API = make(docsonnet.Fields)

	for k, v := range d {
		if k == "#" || !strings.HasPrefix(k, "#") {
			continue
		}

		name := strings.TrimPrefix(k, "#")
		pkg.API[name] = loadField(name, v.(map[string]interface{}), d)
	}

	return pkg
}

func loadField(name string, field map[string]interface{}, parent map[string]interface{}) docsonnet.Field {
	if ifn, ok := field["function"]; ok {
		return loadFn(name, ifn.(map[string]interface{}))
	}

	if iobj, ok := field["object"]; ok {
		return loadObj(name, iobj.(map[string]interface{}), parent)
	}

	panic("docsonnet field lacking {function | object}")
}

func loadFn(name string, msi map[string]interface{}) docsonnet.Field {
	fn := docsonnet.Function{
		Name: name,
		Help: msi["help"].(string),
	}
	if args, ok := msi["args"]; ok {
		fn.Args = loadArgs(args.([]interface{}))
	}
	return docsonnet.Field{Function: &fn}
}

func loadArgs(is []interface{}) []docsonnet.Argument {
	args := make([]docsonnet.Argument, len(is))
	for i := range is {
		arg := is[i].(map[string]interface{})
		args[i] = docsonnet.Argument{
			Name:    arg["name"].(string),
			Type:    docsonnet.Type(arg["type"].(string)),
			Default: arg["default"],
		}
	}
	return args
}

func loadObj(name string, msi map[string]interface{}, parent map[string]interface{}) docsonnet.Field {
	obj := docsonnet.Object{
		Name:   name,
		Help:   msi["help"].(string),
		Fields: make(docsonnet.Fields),
	}

	// look for children in same key without #
	var iChilds interface{}
	var ok bool
	if iChilds, ok = parent[name]; !ok {
		fmt.Println("aborting, no", name)
		return docsonnet.Field{Object: &obj}
	}

	childs := iChilds.(map[string]interface{})
	for k, v := range childs {
		if !strings.HasPrefix(k, "#") {
			continue
		}

		name := strings.TrimPrefix(k, "#")
		obj.Fields[name] = loadField(name, v.(map[string]interface{}), msi)
	}

	return docsonnet.Field{Object: &obj}
}

type DS map[string]interface{}

func (d DS) Package() docsonnet.Package {
	pkg := d["#"].(map[string]interface{})
	return docsonnet.Package{
		Help:   pkg["help"].(string),
		Name:   pkg["name"].(string),
		Import: pkg["import"].(string),
	}
}

func eval() ([]byte, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{JPaths: []string{".."}})
	data, err := ioutil.ReadFile("fast.libsonnet")
	if err != nil {
		return nil, err
	}

	out, err := vm.EvaluateSnippet("fast.libsonnet", string(data))
	if err != nil {
		return nil, err
	}
	return []byte(out), nil
}

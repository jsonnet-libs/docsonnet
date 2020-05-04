package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/google/go-jsonnet"
	"github.com/sh0rez/docsonnet/pkg/docsonnet"
)

func main() {
	data, err := eval()
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(data))

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

// load docsonnet
//
// Data assumptions:
// - only map[string]interface{} and docsonnet fields
// - docsonnet fields (#...) coming first
func load(d DS) docsonnet.Package {
	start := time.Now()

	pkg := d.Package()
	pkg.API = make(docsonnet.Fields)
	pkg.Sub = make(map[string]docsonnet.Package)

	for k, v := range d {
		if k == "#" {
			continue
		}

		f := v.(map[string]interface{})

		// docsonnet field
		name := strings.TrimPrefix(k, "#")
		if strings.HasPrefix(k, "#") {
			pkg.API[name] = loadField(name, f, d)
			continue
		}

		// non-docsonnet
		// subpackage?
		if _, ok := f["#"]; ok {
			p := load(DS(f))
			pkg.Sub[p.Name] = p
			continue
		}

		// non-annotated nested?
		// try to load, but skip when already loaded as annotated above
		if nested, ok := loadNested(name, f); ok && !fieldsHas(pkg.API, name) {
			pkg.API[name] = *nested
			continue
		}
	}

	fmt.Println("load", time.Since(start))
	return pkg
}

func fieldsHas(f docsonnet.Fields, key string) bool {
	_, b := f[key]
	return b
}

func loadNested(name string, msi map[string]interface{}) (*docsonnet.Field, bool) {
	out := docsonnet.Object{
		Name:   name,
		Fields: make(docsonnet.Fields),
	}

	ok := false
	for k, v := range msi {
		f := v.(map[string]interface{})
		n := strings.TrimPrefix(k, "#")

		if !strings.HasPrefix(k, "#") {
			if l, ok := loadNested(k, f); ok {
				out.Fields[n] = *l
			}
			continue
		}

		ok = true
		l := loadField(n, f, msi)
		out.Fields[n] = l
	}

	if !ok {
		return nil, false
	}

	return &docsonnet.Field{Object: &out}, true
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
	hash, ok := d["#"]
	if !ok {
		log.Fatalln("Package declaration missing")
	}

	pkg := hash.(map[string]interface{})
	return docsonnet.Package{
		Help:   pkg["help"].(string),
		Name:   pkg["name"].(string),
		Import: pkg["import"].(string),
	}
}

func eval() ([]byte, error) {
	start := time.Now()

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

	fmt.Println("eval:", time.Since(start))
	return []byte(out), nil
}

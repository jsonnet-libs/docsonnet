package docsonnet

import (
	"fmt"
	"log"
	"strings"
)

// load docsonnet
//
// Data assumptions:
// - only map[string]interface{} and fields
// - fields (#...) coming first
func fastLoad(d ds) Package {
	pkg := d.Package()

	pkg.API = make(Fields)
	pkg.Sub = make(map[string]Package)

	for k, v := range d {
		if k == "#" {
			continue
		}

		n := strings.TrimPrefix(k, "#")
		f := v.(map[string]interface{})

		// is it a docstring?
		if strings.HasPrefix(k, "#") {
			pkg.API[n] = loadField(n, f, d)
			continue
		}

		// is it a package?
		if _, ok := f["#"]; ok {
			p := fastLoad(ds(f))
			pkg.Sub[p.Name] = p
			continue
		}

		// is it a regular field? check nested...
		if nested, ok := loadNested(n, f); ok && !hasDocstring(n, d) {
			pkg.API[n] = *nested
		}
	}

	return pkg
}

func hasDocstring(key string, msi map[string]interface{}) bool {
	_, ok := msi["#"+key]
	return ok
}

func loadNested(name string, msi map[string]interface{}) (*Field, bool) {
	out := Object{
		Name:   name,
		Fields: make(Fields),
	}

	for k, v := range msi {
		n := strings.TrimPrefix(k, "#")
		f := v.(map[string]interface{})

		// is it a docstring?
		if strings.HasPrefix(k, "#") {
			out.Fields[n] = loadField(n, f, msi)
			continue
		}

		// is it a regular field? check nested...
		if nested, ok := loadNested(n, f); ok && !hasDocstring(n, msi) {
			out.Fields[n] = *nested
		}
	}

	return &Field{Object: &out}, true
}

func loadField(name string, field map[string]interface{}, parent map[string]interface{}) Field {
	if ifn, ok := field["function"]; ok {
		return loadFn(name, ifn.(map[string]interface{}))
	}

	if iobj, ok := field["object"]; ok {
		return loadObj(name, iobj.(map[string]interface{}), parent)
	}

	if vobj, ok := field["value"]; ok {
		return loadValue(name, vobj.(map[string]interface{}))
	}

	panic(fmt.Sprintf("field %s lacking {function | object | value}", name))
}

func loadValue(name string, msi map[string]interface{}) Field {
	h, ok := msi["help"].(string)
	if !ok {
		h = ""
	}

	t, ok := msi["type"].(string)
	if !ok {
		panic(fmt.Sprintf("value %s lacking type information", name))
	}

	v := Value{
		Name:    name,
		Help:    h,
		Type:    Type(t),
		Default: msi["default"],
	}

	return Field{Value: &v}
}

func loadFn(name string, msi map[string]interface{}) Field {
	h, ok := msi["help"].(string)
	if !ok {
		h = ""
	}
	fn := Function{
		Name: name,
		Help: h,
	}
	if args, ok := msi["args"]; ok {
		fn.Args = loadArgs(args.([]interface{}))
	}
	return Field{Function: &fn}
}

func loadArgs(is []interface{}) []Argument {
	args := make([]Argument, len(is))
	for i := range is {
		arg := is[i].(map[string]interface{})
		args[i] = Argument{
			Name:    arg["name"].(string),
			Type:    Type(arg["type"].(string)),
			Default: arg["default"],
		}
	}
	return args
}

func fieldNames(msi map[string]interface{}) []string {
	out := make([]string, 0, len(msi))
	for k := range msi {
		out = append(out, k)
	}
	return out
}

func loadObj(name string, msi map[string]interface{}, parent map[string]interface{}) Field {
	obj := Object{
		Name:   name,
		Help:   msi["help"].(string),
		Fields: make(Fields),
	}

	// look for children in same key
	var iChilds interface{}
	var ok bool
	if iChilds, ok = parent[name]; !ok {
		fmt.Println("aborting, no", name, strings.Join(fieldNames(parent), ", "))
		return Field{Object: &obj}
	}

	childs := iChilds.(map[string]interface{})
	if nested, ok := loadNested(name, childs); ok {
		obj.Fields = nested.Object.Fields
	}

	return Field{Object: &obj}
}

type ds map[string]interface{}

func (d ds) Package() Package {
	hash, ok := d["#"]
	if !ok {
		log.Fatalln("Package declaration missing")
	}

	pkg := hash.(map[string]interface{})
	return Package{
		Help:   pkg["help"].(string),
		Name:   pkg["name"].(string),
		Import: pkg["import"].(string),
	}
}

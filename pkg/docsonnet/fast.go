package docsonnet

import (
	"fmt"
	"log"
	"strings"
)

type LoadedField struct {
	DocField        Field
	HasNestedFields bool
	NestedFields    *Field
}

type LoadedFields map[string]*LoadedField

// load docsonnet
//
// Data assumptions:
// - only map[string]interface{} and fields
// - fields (#...) coming first
func fastLoad(d ds) Package {
	pkg := d.Package()

	pkg.API = make(Fields)
	pkg.Sub = make(map[string]Package)

	loadedFields := make(LoadedFields)

	for k, v := range d {
		if k == "#" {
			continue
		}

		n := strings.TrimPrefix(k, "#")
		f := v.(map[string]interface{})

		// is it a package?
		if _, ok := f["#"]; ok {
			p := fastLoad(ds(f))
			pkg.Sub[p.Name] = p
			continue
		}

		// initialize FieldDoc
		if _, ok := loadedFields[n]; !ok {
			loadedFields[n] = &LoadedField{}
		}

		// is it a docstring?
		if strings.HasPrefix(k, "#") {
			loadedFields[n].DocField = loadField(n, f, d)
			continue
		}

		// is it a regular field? check children...
		if nested, ok := loadNested(n, f); ok {
			loadedFields[n].HasNestedFields = true
			loadedFields[n].NestedFields = nested
		}
	}

	pkg.API = consolidateLoadedFields(loadedFields)

	return pkg
}

func loadNested(name string, msi map[string]interface{}) (*Field, bool) {
	out := Object{
		Name:   name,
		Fields: make(Fields),
	}

	loadedFields := make(LoadedFields)

	for k, v := range msi {
		n := strings.TrimPrefix(k, "#")
		f := v.(map[string]interface{})

		// initialize FieldDoc
		if _, ok := loadedFields[n]; !ok {
			loadedFields[n] = &LoadedField{}
		}

		// is it a docstring?
		if strings.HasPrefix(k, "#") {
			loadedFields[n].DocField = loadField(n, f, msi)
			continue
		}

		// is it a regular field? check children...
		if nested, ok := loadNested(n, f); ok {
			loadedFields[n].HasNestedFields = true
			loadedFields[n].NestedFields = nested
		}
	}

	out.Fields = consolidateLoadedFields(loadedFields)

	return &Field{Object: &out}, true
}

func consolidateLoadedFields(loadedFields LoadedFields) Fields {
	fields := make(Fields)
	for k, v := range loadedFields {
		// non-annoted but has childs, only add childs
		if v.DocField == (Field{}) && v.HasNestedFields {
			fields[k] = *v.NestedFields
			continue
		}

		// annotated and has chidls, add childs to docfield
		if v.DocField.Object != nil && v.HasNestedFields {
			v.DocField.Object.Name = k
			for name, child := range v.NestedFields.Object.Fields {
				v.DocField.Object.Fields[name] = child
			}
		}

		// return non-empty docfield
		if v.DocField != (Field{}) {
			fields[k] = v.DocField
		}
	}
	return fields
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

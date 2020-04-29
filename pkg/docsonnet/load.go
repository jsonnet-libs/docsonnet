package docsonnet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/go-jsonnet"
	"github.com/markbates/pkger"
)

// Load extracts docsonnet data from the given Jsonnet document
func Load(filename string) (*Package, error) {
	// get load.libsonnet from embedded data
	file, err := pkger.Open("/load.libsonnet")
	if err != nil {
		return nil, err
	}
	load, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// setup Jsonnet vm
	vm := jsonnet.MakeVM()
	importer, err := newImporter()
	if err != nil {
		return nil, err
	}
	vm.Importer(importer)

	// invoke load.libsonnet
	vm.ExtCode("main", fmt.Sprintf(`(import "%s")`, filename))
	data, err := vm.EvaluateSnippet("load.libsonnet", string(load))
	if err != nil {
		return nil, err
	}

	// parse the result
	var d Package
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatalln(err)
	}

	return &d, nil
}

// importer wraps jsonnet.FileImporter, to statically provide load.libsonnet,
// bundled with the binary
type importer struct {
	fi   jsonnet.FileImporter
	util jsonnet.Contents
}

func newImporter() (*importer, error) {
	file, err := pkger.Open("/doc-util/main.libsonnet")
	if err != nil {
		return nil, err
	}
	load, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &importer{
		fi:   jsonnet.FileImporter{},
		util: jsonnet.MakeContents(string(load)),
	}, nil
}

func (i *importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	if importedPath == "doc-util/main.libsonnet" {
		return i.util, "<internal>", nil
	}

	return i.fi.Import(importedFrom, importedPath)
}

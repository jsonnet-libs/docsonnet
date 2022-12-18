package docsonnet

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/go-jsonnet"
)

type Opts struct {
	JPath      []string
	EmbeddedFS embed.FS
}

// Load extracts and transforms the docsonnet data in `filename`, returning the
// top level docsonnet package.
func Load(filename string, opts Opts) (*Package, error) {
	data, err := Extract(filename, opts)
	if err != nil {
		return nil, err
	}

	return Transform([]byte(data))
}

// Extract parses the Jsonnet file at `filename`, extracting all docsonnet related
// information, exactly as they appear in Jsonnet. Keep in mind this
// representation is usually not suitable for any use, use `Transform` to
// convert it to the familiar docsonnet data model.
func Extract(filename string, opts Opts) ([]byte, error) {
	// get load.libsonnet from embedded data
	load, err := opts.EmbeddedFS.ReadFile("load.libsonnet")
	if err != nil {
		return nil, err
	}

	// setup Jsonnet vm
	vm := jsonnet.MakeVM()
	importer, err := newImporter(opts)
	if err != nil {
		return nil, err
	}
	vm.Importer(importer)

	// invoke load.libsonnet
	vm.ExtCode("main", fmt.Sprintf(`(import "%s")`, filename))

	data, err := vm.EvaluateAnonymousSnippet("load.libsonnet", string(load))
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

// Transform converts the raw result of `Extract` to the actual docsonnet object
// model `*docsonnet.Package`.
func Transform(data []byte) (*Package, error) {
	var d ds
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatalln(err)
	}

	p := fastLoad(d)
	return &p, nil
}

// importer wraps jsonnet.FileImporter, to statically provide load.libsonnet,
// bundled with the binary
type importer struct {
	fi   jsonnet.FileImporter
	util jsonnet.Contents
}

func newImporter(opts Opts) (*importer, error) {
	load, err := opts.EmbeddedFS.ReadFile("doc-util/main.libsonnet")
	if err != nil {
		return nil, err
	}

	return &importer{
		fi:   jsonnet.FileImporter{JPaths: opts.JPath},
		util: jsonnet.MakeContents(string(load)),
	}, nil
}

var docUtilPaths = []string{
	"doc-util/main.libsonnet",
	"github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet",
}

func (i *importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	for _, p := range docUtilPaths {
		if importedPath == p {
			return i.util, "<internal>", nil
		}
	}

	return i.fi.Import(importedFrom, importedPath)
}

package docsonnet

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
)

type Opts struct {
	JPath      []string
	EmbeddedFS embed.FS
}

// RenderWithJsonnet uses the jsonnet render function to generate the docs, instead of the golang utilities.
func RenderWithJsonnet(filename string, opts Opts) (map[string]string, error) {
	// Write out the embedded doc-util to a tmp dir so that we can import it using the native jsonnet importer.
	tmpdir, err := os.MkdirTemp("", "docsonnet-doc-util-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	if err := writeDocUtil(opts.EmbeddedFS, tmpdir); err != nil {
		return nil, err
	}

	// get render.libsonnet from embedded data
	render, err := opts.EmbeddedFS.ReadFile("render.libsonnet")
	if err != nil {
		return nil, err
	}

	// setup Jsonnet vm
	jpaths := append(opts.JPath, tmpdir)
	vm := newVM(filename, jpaths)

	// invoke render.libsonnet
	vm.ExtCode("d", `(import "doc-util/main.libsonnet")`)

	data, err := vm.EvaluateAnonymousSnippet("render.libsonnet", string(render))
	if err != nil {
		return nil, err
	}

	var out map[string]string
	err = json.Unmarshal([]byte(data), &out)
	return out, err
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
	// Write out the embedded doc-util to a tmp dir so that we can import it using the native jsonnet importer.
	tmpdir, err := os.MkdirTemp("", "docsonnet-doc-util-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpdir)

	if err := writeDocUtil(opts.EmbeddedFS, tmpdir); err != nil {
		return nil, err
	}

	// get load.libsonnet from embedded data
	load, err := opts.EmbeddedFS.ReadFile("load.libsonnet")
	if err != nil {
		return nil, err
	}

	// setup Jsonnet vm
	jpaths := append(opts.JPath, tmpdir)
	vm := newVM(filename, jpaths)

	// invoke load.libsonnet
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

// newVM sets up the Jsonnet VM with the importer that statically provides doc-util.
func newVM(mainFName string, jpaths []string) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{JPaths: jpaths})
	vm.ExtCode("main", fmt.Sprintf(`(import "%s")`, mainFName))
	return vm
}

// writeDocUtil writes the embedded doc-util libsonnet package to disk so that Jsonnet can load it.
func writeDocUtil(embedded embed.FS, tmpdir string) error {
	rootDir := filepath.Join(tmpdir, "github.com/jsonnet-libs/docsonnet/doc-util")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}
	spath := filepath.Join(tmpdir, "doc-util")
	if err := os.Symlink(rootDir, spath); err != nil {
		return err
	}

	dir, err := embedded.ReadDir("doc-util")
	if err != nil {
		return err
	}
	for _, f := range dir {
		fpath := filepath.Join("doc-util", f.Name())
		conts, err := embedded.ReadFile(fpath)
		if err != nil {
			return err
		}

		outPath := filepath.Join(rootDir, f.Name())
		if err := os.WriteFile(outPath, conts, 0644); err != nil {
			return err
		}
	}
	return nil
}

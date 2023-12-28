package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-clix/cli"
	"github.com/google/go-jsonnet"
)

var (
	//go:embed doc-util
	embedded embed.FS
)

func main() {
	log.SetFlags(0)

	root := &cli.Command{
		Use:   "docsonnet <file>",
		Short: "Utility to parse and transform Jsonnet code that uses the docsonnet extension",
		Args:  cli.ArgsExact(1),
	}

	dir := root.Flags().StringP("output", "o", "docs", "directory to write the .md files to")
	jpath := root.Flags().StringSliceP("jpath", "J", []string{"vendor"}, "Specify an additional library search dir (right-most wins)")

	root.Run = func(cmd *cli.Command, args []string) error {
		file := args[0]

		vm := jsonnet.MakeVM()
		importer, err := newImporter(*jpath)
		if err != nil {
			return err
		}
		vm.Importer(importer)

		renderSnippet, err := embedded.ReadFile("doc-util/render.libsonnet")
		if err != nil {
			return err
		}

		snippet := fmt.Sprintf(`(%s).render(std.extVar('main'))`, renderSnippet)

		vm.ExtCode("main", fmt.Sprintf(`(import "%s")`, file))

		jsonStr, err := vm.EvaluateAnonymousSnippet("main.libsonnet", snippet)
		if err != nil {
			return err
		}

		var output map[string]string
		err = json.Unmarshal([]byte(jsonStr), &output)
		if err != nil {
			return err
		}

		for k, v := range output {
			fullpath := filepath.Join(*dir, k)
			dir := filepath.Dir(fullpath)
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
			if err := os.WriteFile(fullpath, []byte(v), 0644); err != nil {
				return err
			}
		}
		return nil
	}

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

// importer wraps jsonnet.FileImporter
type importer struct {
	fi   jsonnet.FileImporter
	util jsonnet.Contents
}

func newImporter(paths []string) (*importer, error) {
	load, err := embedded.ReadFile("doc-util/main.libsonnet")
	if err != nil {
		return nil, err
	}

	render, err := embedded.ReadFile("doc-util/render.libsonnet")
	if err != nil {
		return nil, err
	}

	main := strings.ReplaceAll(string(load), "(import './render.libsonnet')", string(render))

	return &importer{
		fi:   jsonnet.FileImporter{JPaths: paths},
		util: jsonnet.MakeContents(main),
	}, nil
}

var docUtilPaths = []string{
	"doc-util/main.libsonnet",
	"github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet",
}

func (i *importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	for _, p := range docUtilPaths {
		if importedPath == p {
			return i.util, "main.libsonnet", nil
		}
	}

	return i.fi.Import(importedFrom, importedPath)
}

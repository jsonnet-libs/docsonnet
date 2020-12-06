package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-clix/cli"

	"github.com/sh0rez/docsonnet/pkg/docsonnet"
	"github.com/sh0rez/docsonnet/pkg/render"
)

var Version = "dev"

func main() {
	log.SetFlags(0)

	root := &cli.Command{
		Use:     "docsonnet <file>",
		Short:   "Utility to parse and transform Jsonnet code that uses the docsonnet extension",
		Args:    cli.ArgsExact(1),
		Version: Version,
	}

	dir := root.Flags().StringP("output", "o", "docs", "directory to write the .md files to")
	outputJSON := root.Flags().Bool("json", false, "print loaded docsonnet as JSON")
	outputRaw := root.Flags().Bool("raw", false, "don't transform, dump raw eval result")
	urlPrefix := root.Flags().String("urlPrefix", "/", "url-prefix for frontmatter")
	jpath := root.Flags().StringSliceP("jpath", "J", []string{"vendor"}, "Specify an additional library search dir (right-most wins)")

	root.Run = func(cmd *cli.Command, args []string) error {
		file := args[0]

		log.Println("Extracting from Jsonnet")
		data, err := docsonnet.Extract(file, docsonnet.Opts{JPath: *jpath})
		if err != nil {
			log.Fatalln("Extracting:", err)
		}
		if *outputRaw {
			fmt.Println(string(data))
			return nil
		}

		log.Println("Transforming to docsonnet model")
		pkg, err := docsonnet.Transform(data)
		if err != nil {
			log.Fatalln("Transforming:", err)
		}
		if *outputJSON {
			data, err := json.MarshalIndent(pkg, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		log.Println("Rendering markdown")
		n, err := render.To(*pkg, *dir, render.Opts{
			URLPrefix: *urlPrefix,
		})
		if err != nil {
			log.Fatalln("Rendering:", err)
		}

		log.Printf("Success! Rendered %v packages from '%s' to '%s'", n, file, *dir)
		return nil
	}

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

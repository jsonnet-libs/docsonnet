package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-clix/cli"

	"github.com/sh0rez/docsonnet/pkg/docsonnet"
	"github.com/sh0rez/docsonnet/pkg/render"
)

func main() {
	log.SetFlags(0)

	root := &cli.Command{
		Use:   "docsonnet",
		Short: "Utility to parse and transform Jsonnet code that uses the docsonnet extension",
	}

	dir := root.Flags().StringP("output", "o", "docs", "directory to write the .md files to")
	outputMd := root.Flags().Bool("md", true, "render as markdown files")
	outputJSON := root.Flags().Bool("json", false, "print loaded docsonnet as JSON")

	root.Run = func(cmd *cli.Command, args []string) error {
		file := args[0]

		switch {
		case *outputJSON:
			model, err := docsonnet.Load(file)
			if err != nil {
				return err
			}
			data, err := json.MarshalIndent(model, "", "  ")
			if err != nil {
				return err
			}

			fmt.Println(string(data))
		case *outputMd:
			return render.To(file, *dir)
		}

		return nil
	}

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

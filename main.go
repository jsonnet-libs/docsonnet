package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

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

	output := root.Flags().StringP("output", "o", "docs", "directory to write the .md files to")
	root.Run = func(cmd *cli.Command, args []string) error {
		pkg, err := docsonnet.Load(args[0])
		if err != nil {
			return err
		}

		data := render.Render(*pkg)
		for k, v := range data {
			fmt.Println(k)
			if err := ioutil.WriteFile(filepath.Join(*output, k), []byte(v), 0644); err != nil {
				return err
			}
		}

		return nil
	}

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

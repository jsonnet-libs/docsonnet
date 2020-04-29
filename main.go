package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	root.AddCommand(loadCmd(), renderCmd())

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func loadCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "load",
		Short: "extracts docsonnet from Jsonnet and prints it as JSON",
		Args:  cli.ArgsExact(1),
	}

	cmd.Run = func(cmd *cli.Command, args []string) error {
		pkg, err := docsonnet.Load(args[0])
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return cmd
}

func renderCmd() *cli.Command {
	cmd := &cli.Command{
		Use:   "render",
		Short: "writes all found docsonnet packages to Markdown (.md) files, suitable for e.g. GitHub",
		Args:  cli.ArgsExact(1),
	}

	output := cmd.Flags().StringP("output", "o", "docs", "directory to write the .md files to")

	cmd.Run = func(cmd *cli.Command, args []string) error {
		pkg, err := docsonnet.Load(args[0])
		if err != nil {
			return err
		}

		for path, pkg := range render.Paths(*pkg) {
			to := filepath.Join(*output, path)
			if err := os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
				return err
			}

			data := render.Render(pkg)
			if err := ioutil.WriteFile(to, []byte(data), 0644); err != nil {
				return err
			}
		}

		return nil
	}

	return cmd
}

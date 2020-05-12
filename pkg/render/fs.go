package render

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/sh0rez/docsonnet/pkg/docsonnet"
)

func To(api, out string) error {
	pkg, err := docsonnet.Load(api)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		return err
	}

	log.Println("Rendering .md files")
	data := Render(*pkg)
	for k, v := range data {
		if err := ioutil.WriteFile(filepath.Join(out, k), []byte(v), 0644); err != nil {
			return err
		}
	}

	log.Printf("Success! Rendered %v packages from '%s' to '%s'", len(data), api, out)
	return nil
}

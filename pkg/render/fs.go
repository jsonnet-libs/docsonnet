package render

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jsonnet-libs/docsonnet/pkg/docsonnet"
)

func To(pkg docsonnet.Package, dir string, opts Opts) (int, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return 0, err
	}

	data := Render(pkg, opts)

	n := 0
	for k, v := range data {
		fullpath := filepath.Join(dir, k)
		dir := filepath.Dir(fullpath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return n, err
		}
		if err := ioutil.WriteFile(fullpath, []byte(v), 0644); err != nil {
			return n, err
		}
		n++
	}

	return n, nil
}

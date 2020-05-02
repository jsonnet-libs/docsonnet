package render

import (
	"github.com/sh0rez/docsonnet/pkg/docsonnet"
)

func Paths(pkg docsonnet.Package) map[string]docsonnet.Package {
	p := paths(pkg)
	return p
}

func paths(pkg docsonnet.Package) map[string]docsonnet.Package {
	pkgs := make(map[string]docsonnet.Package)
	pkgs[pkg.Name+".md"] = pkg

	if len(pkg.Sub) == 0 {
		return pkgs
	}

	for _, sub := range pkg.Sub {
		for k, v := range paths(sub) {
			v.Name = pkg.Name + "/" + k
			pkgs[pkg.Name+"-"+k] = v
		}
	}

	return pkgs
}

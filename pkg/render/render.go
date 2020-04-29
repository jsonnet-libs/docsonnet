package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sh0rez/docsonnet/pkg/docsonnet"
	"github.com/sh0rez/docsonnet/pkg/md"
	"github.com/sh0rez/docsonnet/pkg/slug"
)

func Render(pkg docsonnet.Package) string {
	// head
	elems := []md.Elem{
		md.Headline(1, "package "+pkg.Name),
		md.CodeBlock("jsonnet", fmt.Sprintf(`local %s = import "%s"`, pkg.Name, pkg.Import)),
		md.Text(pkg.Help),
	}

	// subpackages
	if len(pkg.Sub) > 0 {
		elems = append(elems, md.Headline(2, "Subpackages"))
		var items []md.Elem
		for k := range Paths(pkg) {
			name := strings.TrimSuffix(k, ".md")
			if name == pkg.Name {
				continue
			}
			items = append(items, md.Link(md.Text(name), k))
		}
		elems = append(elems, md.List(items...))
	}

	// fields of this package
	if len(pkg.API) > 0 {
		// index
		elems = append(elems,
			md.Headline(2, "Index"),
			md.List(renderIndex(pkg.API, "", slug.New())...),
		)

		// api
		elems = append(elems, md.Headline(2, "Fields"))
		elems = append(elems, renderApi(pkg.API, "")...)
	}

	return md.Doc(elems...).String()
}

func renderIndex(api docsonnet.Fields, path string, s *slug.Slugger) []md.Elem {
	var elems []md.Elem
	for _, k := range sortFields(api) {
		v := api[k]
		switch {
		case v.Function != nil:
			fn := v.Function
			name := md.Text(fmt.Sprintf("fn %s(%s)", fn.Name, renderParams(fn.Args)))
			link := "#" + s.Slug("fn "+path+fn.Name)
			elems = append(elems, md.Link(md.Code(name), link))
		case v.Object != nil:
			obj := v.Object
			name := md.Text("obj " + path + obj.Name)
			link := "#" + s.Slug("obj "+path+obj.Name)
			elems = append(elems, md.Link(md.Code(name), link))
			elems = append(elems, md.List(renderIndex(obj.Fields, path+obj.Name+".", s)...))
		}
	}
	return elems
}

func renderApi(api docsonnet.Fields, path string) []md.Elem {
	var elems []md.Elem

	for _, k := range sortFields(api) {
		v := api[k]
		switch {
		case v.Function != nil:
			fn := v.Function
			elems = append(elems,
				md.Headline(3, fmt.Sprintf("fn %s%s", path, fn.Name)),
				md.CodeBlock("ts", fmt.Sprintf("%s(%s)", fn.Name, renderParams(fn.Args))),
				md.Text(fn.Help),
			)
		case v.Object != nil:
			obj := v.Object
			elems = append(elems,
				md.Headline(2, fmt.Sprintf("obj %s%s", path, obj.Name)),
				md.Text(obj.Help),
			)
			elems = append(elems, renderApi(obj.Fields, path+obj.Name+".")...)
		}
	}

	return elems
}

func sortFields(api docsonnet.Fields) []string {
	keys := make([]string, len(api))
	for k := range api {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		iK, jK := keys[i], keys[j]
		if api[iK].Function != nil && api[jK].Function == nil {
			return true
		}
		if api[iK].Function == nil && api[jK].Function != nil {
			return false
		}
		return iK < jK
	})

	return keys
}

func renderParams(a []docsonnet.Argument) string {
	args := make([]string, 0, len(a))
	for _, a := range a {
		arg := a.Name
		if a.Default != nil {
			arg = fmt.Sprintf("%s=%v", arg, a.Default)
		}
		args = append(args, arg)
	}

	return strings.Join(args, ", ")
}

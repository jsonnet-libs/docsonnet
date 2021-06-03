package render

import (
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/google/go-jsonnet/formatter"
	"github.com/jsonnet-libs/docsonnet/pkg/docsonnet"
	"github.com/jsonnet-libs/docsonnet/pkg/md"
	"github.com/jsonnet-libs/docsonnet/pkg/slug"
)

type Opts struct {
	URLPrefix string
}

func Render(pkg docsonnet.Package, opts Opts) map[string]string {
	return render(pkg, nil, true, opts.URLPrefix)
}

func render(pkg docsonnet.Package, parents []string, root bool, urlPrefix string) map[string]string {
	link := path.Join("/", urlPrefix, strings.Join(append(parents, pkg.Name), "/"))
	if root {
		link = path.Join("/", urlPrefix)
	}
	if !strings.HasSuffix(link, "/") {
		link = link + "/"
	}

	// head
	elems := []md.Elem{
		md.Frontmatter(map[string]interface{}{
			"permalink": link,
		}),
		md.Headline(1, "package "+pkg.Name),
	}
	if pkg.Import != "" {
		elems = append(elems, md.CodeBlock("jsonnet", fmt.Sprintf(`local %s = import "%s"`, pkg.Name, pkg.Import)))
	}
	elems = append(elems, md.Text(pkg.Help))

	if len(pkg.Sub) > 0 {
		elems = append(elems, md.Headline(2, "Subpackages"))

		keys := make([]string, 0, len(pkg.Sub))
		for _, s := range pkg.Sub {
			keys = append(keys, s.Name)
		}
		sort.Strings(keys)

		var items []md.Elem
		for _, k := range keys {
			s := pkg.Sub[k]

			link := strings.Join(append(parents, pkg.Name, s.Name), "-")
			if root {
				link = strings.Join(append(parents, s.Name), "-")
			}
			items = append(items, md.Link(md.Text(s.Name), link+".md"))
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

	content := md.Doc(elems...).String()
	key := strings.Join(append(parents, pkg.Name+".md"), "-")
	if root {
		key = "README.md"
	}
	out := map[string]string{
		key: content,
	}

	if len(pkg.Sub) != 0 {
		for _, s := range pkg.Sub {
			path := append(parents, pkg.Name)
			if root {
				path = parents
			}
			got := render(s, path, false, urlPrefix)
			for k, v := range got {
				out[k] = v
			}
		}
	}

	return out
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
		case v.Value != nil:
			val := v.Value
			name := md.Text(fmt.Sprintf("%s %s%s", val.Type, path, val.Name))
			link := "#" + s.Slug(name.String())
			elems = append(elems, md.Link(md.Code(name), link))
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

		case v.Value != nil:
			val := v.Value
			elems = append(elems,
				md.Headline(3, fmt.Sprintf("%s %s%s", val.Type, path, val.Name)),
			)

			if val.Default != nil {
				elems = append(elems, md.Paragraph(
					md.Italic(md.Text("Default value: ")),
					md.Code(md.Text(fmt.Sprint(val.Default))),
				))
			}

			elems = append(elems,
				md.Text(val.Help),
			)
		}
	}

	return elems
}

func sortFields(api docsonnet.Fields) []string {
	keys := make([]string, 0, len(api))
	for k := range api {
		keys = append(keys, k)
	}

	aFn := func(a, b string) bool {
		return api[a].Function != nil && api[b].Function == nil
	}
	aNew := func(a, b string) bool {
		a = strings.ToLower(a)
		b = strings.ToLower(b)

		return strings.HasPrefix(a, "new") && !strings.HasPrefix(b, "new")
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]

		if aNew(a, b) {
			return true
		} else if aNew(b, a) {
			return false
		}

		if aFn(a, b) {
			return true
		} else if aFn(b, a) {
			return false
		}

		return a < b
	})

	return keys
}

func renderParams(a []docsonnet.Argument) string {
	args := make([]string, 0, len(a))
	for _, a := range a {
		arg := a.Name
		if a.Default != nil {
			arg = fmt.Sprintf("%s=%v", arg, jsonParam(a.Default))
		}
		args = append(args, arg)
	}

	return strings.Join(args, ", ")
}

func jsonParam(i interface{}) string {

	d, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	s, err := formatter.Format("(jsonParam)", string(d), formatter.Options{
		PadObjects:       false,
		PadArrays:        false,
		PrettyFieldNames: true,
		StringStyle:      formatter.StringStyleSingle,
	})
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(s)
}

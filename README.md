# `docsonnet`

This repository contains an experimental Jsonnet docs generator, consisting of multiple parts:

- **Docsonnet**, a **data model** for logically describing the structure of public
  facing Jsonnet API's.
- `doc-util`, A **Jsonnet extension** that allows to write Docsonnet directly
  alongside your Jsonnet. Currently implemented as a library, might become
  language sugar at some point
- `docsonnet`: A **CLI application** and Go library for parsing Docsonnet and
  transforming it to e.g. **Markdown** pages

## Example

To make use of Docsonnet, use `doc-util` to annotate your Jsonnet like so:

```jsonnet
{
    // package declaration
    '#': d.pkg(
      name='url',
      url='github.com/jsonnet-libs/xtd/url/main.libsonet',
      help='`url` implements URL escaping and query building',
    ),

    // function description
    '#encodeQuery': d.fn(
      '`encodeQuery` takes an query parameters and returns them as an escaped `key=value` string',
      [d.arg('params', d.T.object)]),
    encodeQuery(params)::
      local fmtParam(p) = '%s=%s' % [self.escapeString(p), self.escapeString(params[p])];
      std.join('&', std.map(fmtParam, std.objectFields(params))),
}
```

### Packages

Jsonnet itself does not know traditional packages, classes or similar.

For documentation and distribution purposes however, it seems reasonable to introduce a concept of **loose packages**, defined as a single importable file, holding all of your **public API**.

As an example, a hypothetical `url` library could define its package like above example does.

Packages are defined by including assigning a `d.pkg` call to a key literally named `#` (hash). All fields, including nested packages, of the same object having the `#` key belong to that package.

### Functions

Most common part of an API will be functions. These are annotated in a similar fashion:

```jsonnet
{
    "#myFunc": d.fn("myFunc greets you", [d.arg("who", d.T.string)])
    myFunc(who):: "hello %s!" % who
}
```

Along the actual function definition, a _docsonnet_ key is added, with the functions name prefixed by the familiar `#` as its name.
Above example defines `myFunc` as a function, that greets the user and takes a single argument of type `string`.

### Objects

Sometimes you might want to group functions of a similar kind, by nesting them into plain Jsonnet objects.

Such an object might need a description as well, so you can also annotate it:

```jsonnet
{
    "#myObj": d.obj("myObj holds my functions")
    myObj:: {
        "#myFunc": d.fn("myFunc greets you", [d.arg("who", d.T.string)])
        myFunc(who):: "hello %s!" % who
    }
}
```

Again, the naming rule `#` joined with the fields name must be followed, so the `docsonnet` utility can automatically join together the contents of your object with its annotated description.


## Usage

Once you have a Jsonnet library annotated with `doc-util`, you can generate the docs using one of three ways:

- [Jsonnet renderer](#jsonnet-renderer)
- [docsonnet binary](#docsonnet-binary)
- [docsonnet docker image](#docsonnet-docker-image)

### Jsonnet renderer

The docs can be rendered using Jsonnet with the
[render](https://github.com/jsonnet-libs/docsonnet/tree/master/doc-util#fn-render) function.

In your library source, add a file `docs.jsonnet` (assuming your library entrypoint is `main.libsonnet`) with the
following contents:

```jsonnet
local d = import 'github.com/jsonnet-libs/docsonnet/doc-util/main.libsonnet';
d.render(import 'main.libsonnet')
```

Then, you can render the markdown to the `docs/` folder using the following command:

```
jsonnet -S -c -m docs/ docs.jsonnet
```

Note that this requires `doc-util` to be installed to `vender/` to work properly.

### docsonnet binary

Alternatively, the docs can be rendered using the `docsonnet` go binary. The `docsonnet` binary embeds the `doc-util`
library, avoiding the need to have `doc-util` installed.

You can install the `docsonnet` binary using `go install`:

```
go install github.com/jsonnet-libs/docsonnet@master
```

Once the binary is installed, you can generate the docs by passing it the main entrypoint to your Jsonnet library:

```
docsonnet main.libsonnet
```

> **Note**
>
> Linters like [jsonnet-lint](https://pkg.go.dev/github.com/google/go-jsonnet/linter) or `tk lint` require the imports to be resolvable, so you should add `doc-util` to `vendor/` when using these linters.

### docsonnet docker image

You can also use the [docker image](https://hub.docker.com/r/jsonnetlibs/docsonnet) which contains the `docsonnet`
binary if you do not wish to set up go or install the binary locally:

```
docker run --rm -v "$(pwd):/src" -v "$(pwd)/docs:/docs" jsonnetlibs/docsonnet /src/main.libsonnet
```


## FAQ

#### What's wrong with comments? Why not parse regular comments?

I had some attempts on this, especially because it feels more natural. However, the language properties of Jsonnet make this quite challenging:

- AST parsing is insufficient:
  https://github.com/grafana/tanka/issues/223#issuecomment-590569198. Just by
  parsing the syntax tree of Jsonnet, we only receive a representation of the
  file contents, not the logical ones a human might infer
- No effective view on things: Jsonnet is a lazily evaluated, highly dynamic
  language. Just by looking at a single file, we might not even see what ends up
  at the user when importing the library, because during evaluation things can
  be heavily overwritten.

Because of that, we would need to perform a slimmed down evaluation on the AST before getting our information out of it. This is a lot of work, especially when we can just use the real Jsonnet compiler to do this for us. That's docsonnet.

#### But docsonnet is ugly. And verbose

I know. Think of docsonnet as a proof of concept and a technology preview. Only _what_ you specify is a fixed thing, not the way you do.

Of course nobody wants these ugly function calls as docs. But they are incredibly powerful, because we can use Jsonnet merging and patching on the generated docsonnet fields, and the Jsonnet compiler handles that for us.

In case this idea works out well, we might very well consider adding docsonnet as language sugar to Jsonnet, which might look like this:

```jsonnet
{
    ## myFunc greets you
    ## @params:
    ##   who: string
    myFunc(who):: "hello %s!" % who
}
```

Note the double hash `##` as a special indicator for the compiler, so it can desugar above to:

```jsonnet
{
    "#myFunc": d.fn("myFunc greets you", [d.arg("who", d.T.string)])
    myFunc(who):: "hello %s!" % who
}
```

This will all happen transparently, without any user interaction

#### What else can it do?

Because the Docsonnet gives you the missing logical representation of your Jsonnet library, it enables straight forward implementation of other language tooling, such as **code-completion**.

Instead of inferring what fields are available for a library, we can _just_ look at its docsonnet and provide the fields specified there, along with nice descriptions and argument types.

# package d

`doc-util` provides a Jsonnet interface for `docsonnet`,
 a Jsonnet API doc generator that uses structured data instead of comments.


## Install

```
jb install github.com/jsonnet-libs/docsonnet/doc-util@master
```

## Usage

```jsonnet
local d = import "github.com/jsonnet-libs/docsonnet/doc-util/doc-util/main.libsonnet"
```

## Index

* [`fn arg(name, type, default)`](#fn-arg)
* [`fn fn(help, args)`](#fn-fn)
* [`fn obj(help, fields)`](#fn-obj)
* [`fn pkg(name, url, help, filename='', version='master')`](#fn-pkg)
* [`fn render(obj)`](#fn-render)
* [`fn val(type, help, default)`](#fn-val)
* [`obj argument`](#obj-argument)
  * [`fn new(name, type, default)`](#fn-argumentnew)
* [`obj func`](#obj-func)
  * [`fn new(help, args)`](#fn-funcnew)
  * [`fn withArgs(args)`](#fn-funcwithargs)
  * [`fn withHelp(help)`](#fn-funcwithhelp)
* [`obj object`](#obj-object)
  * [`fn new(help, fields)`](#fn-objectnew)
  * [`fn withFields(fields)`](#fn-objectwithfields)
* [`obj value`](#obj-value)
  * [`fn new(type, help, default)`](#fn-valuenew)
* [`obj T`](#obj-t)
* [`obj package`](#obj-package)
  * [`fn new(name, url, help, filename='', version='master')`](#fn-packagenew)

## Fields

### fn arg

```ts
arg(name, type, default)
```

`arg` is a shorthand for `argument.new`

### fn fn

```ts
fn(help, args)
```

`fn` is a shorthand for `func.new`

### fn obj

```ts
obj(help, fields)
```

`obj` is a shorthand for `object.new`

### fn pkg

```ts
pkg(name, url, help, filename='', version='master')
```

`new` is a shorthand for `package.new`

### fn render

```ts
render(obj)
```

`render` converts the docstrings to human readable Markdown files.

Usage:

```jsonnet
// docs.jsonnet
d.render(import 'main.libsonnet')
```

Call with: `jsonnet -S -c -m docs/ docs.jsonnet`


### fn val

```ts
val(type, help, default)
```

`val` is a shorthand for `value.new`

### obj argument

Utilities for creating function arguments

#### fn argument.new

```ts
new(name, type, default)
```

new creates a new function argument, taking the name, the type and optionally a default value

### obj func

Utilities for documenting Jsonnet methods (functions of objects)

#### fn func.new

```ts
new(help, args)
```

new creates a new function, optionally with description and arguments

#### fn func.withArgs

```ts
withArgs(args)
```

The `withArgs` modifier overrides the arguments of that function

#### fn func.withHelp

```ts
withHelp(help)
```

The `withHelp` modifier overrides the help text of that function

### obj object

Utilities for documenting Jsonnet objects (`{ }`).

#### fn object.new

```ts
new(help, fields)
```

new creates a new object, optionally with description and fields

#### fn object.withFields

```ts
withFields(fields)
```

The `withFields` modifier overrides the fields property of an already created object

### obj value

Utilities for documenting plain Jsonnet values (primitives)

#### fn value.new

```ts
new(type, help, default)
```

new creates a new object of given type, optionally with description and default value

### obj T


* `T.any` (`string`): `"any"` - argument of type "any"
* `T.array` (`string`): `"array"` - argument of type "array"
* `T.boolean` (`string`): `"bool"` - argument of type "boolean"
* `T.func` (`string`): `"function"` - argument of type "func"
* `T.null` (`string`): `"null"` - argument of type "null"
* `T.number` (`string`): `"number"` - argument of type "number"
* `T.object` (`string`): `"object"` - argument of type "object"
* `T.string` (`string`): `"string"` - argument of type "string"

### obj package


#### fn package.new

```ts
new(name, url, help, filename='', version='master')
```

`new` creates a new package

Arguments:

* given `name`
* source `url` for jsonnet-bundler and the import
* `help` text
* `filename` for the import, defaults to blank for backward compatibility
* `version` for jsonnet-bundler install, defaults to `master` just like jsonnet-bundler


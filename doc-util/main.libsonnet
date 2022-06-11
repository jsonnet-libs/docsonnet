{
  local d = self,

  '#': d.pkg(
    name='d',
    url='github.com/jsonnet-libs/docsonnet/doc-util',
    help=|||
      `doc-util` provides a Jsonnet interface for `docsonnet`,
       a Jsonnet API doc generator that uses structured data instead of comments.
    |||,
    filename=std.thisFile,
  ),

  package:: {
    '#new':: d.fn(|||
      `new` creates a new package

      Arguments:

      * given `name`
      * source `url` for jsonnet-bundler and the import
      * `help` text
      * `filename` for the import, defaults to blank for backward compatibility
      * `tag` for jsonnet-bundler install, defaults to `master` just like jsonnet-bundler
    |||, [
      d.arg('name', d.T.string),
      d.arg('url', d.T.string),
      d.arg('help', d.T.string),
      d.arg('filename', d.T.string, ''),
      d.arg('tag', d.T.string, 'master'),
    ]),
    new(name, url, help, filename='', tag='master'):: {
      name: name,
      help: help,

      url: url,
      filename: filename,
      tag: tag,

      'import': url + (if filename != '' then '/' + filename else ''),
    },
  },

  '#pkg':: self.package['#new'] + d.func.withHelp('`new` is a shorthand for `package.new`'),
  pkg:: self.package.new,

  '#object': d.obj('Utilities for documenting Jsonnet objects (`{ }`).'),
  object:: {
    '#new': d.fn('new creates a new object, optionally with description and fields', [d.arg('help', d.T.string), d.arg('fields', d.T.object)]),
    new(help='', fields={}):: { object: {
      help: help,
      fields: fields,
    } },

    '#withFields': d.fn('The `withFields` modifier overrides the fields property of an already created object', [d.arg('fields', d.T.object)]),
    withFields(fields):: { object+: {
      fields: fields,
    } },
  },

  '#obj': self.object['#new'] + d.func.withHelp('`obj` is a shorthand for `object.new`'),
  obj:: self.object.new,

  '#func': d.obj('Utilities for documenting Jsonnet methods (functions of objects)'),
  func:: {
    '#new': d.fn('new creates a new function, optionally with description and arguments', [d.arg('help', d.T.string), d.arg('args', d.T.array)]),
    new(help='', args=[]):: { 'function': {
      help: help,
      args: args,
    } },

    '#withHelp': d.fn('The `withHelp` modifier overrides the help text of that function', [d.arg('help', d.T.string)]),
    withHelp(help):: { 'function'+: {
      help: help,
    } },

    '#withArgs': d.fn('The `withArgs` modifier overrides the arguments of that function', [d.arg('args', d.T.array)]),
    withArgs(args):: { 'function'+: {
      args: args,
    } },
  },

  '#fn': self.func['#new'] + d.func.withHelp('`fn` is a shorthand for `func.new`'),
  fn:: self.func.new,

  '#argument': d.obj('Utilities for creating function arguments'),
  argument:: {
    '#new': d.fn('new creates a new function argument, taking the name, the type and optionally a default value', [d.arg('name', d.T.string), d.arg('type', d.T.string), d.arg('default', d.T.any)]),
    new(name, type, default=null): {
      name: name,
      type: type,
      default: default,
    },
  },
  '#arg': self.argument['#new'] + self.func.withHelp('`arg` is a shorthand for `argument.new`'),
  arg:: self.argument.new,

  '#value': d.obj('Utilities for documenting plain Jsonnet values (primitives)'),
  value:: {
    '#new': d.fn('new creates a new object of given type, optionally with description and default value', [d.arg('type', d.T.string), d.arg('help', d.T.string), d.arg('default', d.T.any)]),
    new(type, help='', default=null): { value: {
      help: help,
      type: type,
      default: default,
    } },
  },
  '#val': self.value['#new'] + self.func.withHelp('`val` is a shorthand for `value.new`'),
  val:: self.value.new,

  // T contains constants for the Jsonnet types
  T:: {
    '#string': d.val(d.T.string, 'argument of type "string"'),
    string: 'string',

    '#number': d.val(d.T.string, 'argument of type "number"'),
    number: 'number',
    int: self.number,
    integer: self.number,

    '#boolean': d.val(d.T.string, 'argument of type "boolean"'),
    boolean: 'bool',
    bool: self.boolean,

    '#object': d.val(d.T.string, 'argument of type "object"'),
    object: 'object',

    '#array': d.val(d.T.string, 'argument of type "array"'),
    array: 'array',

    '#any': d.val(d.T.string, 'argument of type "any"'),
    any: 'any',

    '#null': d.val(d.T.string, 'argument of type "null"'),
    'null': 'null',
    nil: self['null'],

    '#func': d.val(d.T.string, 'argument of type "func"'),
    func: 'function',
    'function': self.func,
  },

  '#render': d.fn(
    |||
      `render` converts the docstrings to human readable Markdown files.

      Usage:

      ```jsonnet
      // docs.jsonnet
      d.render(import 'main.libsonnet')
      ```

      Call with: `jsonnet -S -c -m docs/ docs.jsonnet`
    |||,
    args=[
      d.arg('obj', d.T.object),
    ]
  ),
  render:: (import './render.libsonnet').render,

}

local _internal = {
  // clean removes the hash from the key
  clean(str): std.strReplace(str, '#', ''),

  include(api, key):
    // include docsonnet fields
    if std.startsWith(key, '#') then true

    // non-docsonnet object, might have nested docsonnet
    else if std.isObject(api[key]) then
      // skip if has docsonnet counterpart
      (if std.objectHasAll(api, '#' + key) then false else true)

    // some other field, skip it
    else false,

  // ds deals with docsonnet fields. Also checks for nested objects for
  //docsonnet ones
  ds(api, key):
    local cleaned = self.clean(key);
    if std.objectHasAll(api[key], 'object') && std.objectHasAll(api, cleaned) then
      api[key] + $.object.withFields($.filter(api[cleaned]))
    else
      api[key],

  // filter returns only docsonnet related objects from the api
  filter(api): {
    [self.clean(key)]:
      // include all docsonnet
      if std.startsWith(key, '#') then
        $.ds(api, key)
      // create docsonnet objects from regular ones
      else if std.isObject(api[key]) then
        $.obj('', $.filter(api[key]))

    for key in std.objectFields(api)
    if self.include(api, key)
  },

  fromMain(main): self.filter(main),
};

{
  local internal = self + _internal,
  local d = self,

  '#new': d.fn('`new` initiates the api documentation model, taking the `name` and import `url` of your project', [d.arg('name', d.T.string), d.arg('url', d.T.string)]),
  new(name, url):: {
    name: name,
    'import': url,
  },

  '#withAPI': d.fn('`withAPI` automatically builds the docsonnet model from your public API. The `main` parameter is usually set to `import "./main.libsonnet"`.', [d.arg('main', d.T.object)]),
  withAPI(main):: {
    api: internal.fromMain(main),
  },

  '#withFields': d.fn('`withFields` allows to manually specify the docsonnet model. Usually `withAPI` is the better alternative', [d.arg('fields', d.T.object)]),
  withFields(fields):: {
    api: fields,
  },

  "#withHelp": d.fn("`withHelp` sets the main description text, displayed right under the package's name", [d.arg("help", d.T.string)]),
  withHelp(help):: {
    help: help,
  },

  "#object": d.obj("Utilities for documenting Jsonnet objects (`{ }`)."),
  object:: {
    "#new": d.fn("new creates a new object, optionally with description and fields", [d.arg("help", d.T.string), d.arg("fields", d.T.object)]),
    new(help='', fields={}):: { object: {
      help: help,
      fields: fields,
    } },

    "#withFields": d.fn("The `withFields` modifier overrides the fields property of an already created object", [d.arg("fields", d.T.object)]),
    withFields(fields):: { object+: {
      fields: fields,
    } },
  },

  "#obj": self.object["#new"] + d.func.withHelp("`obj` is a shorthand for `object.new`"),
  obj:: self.object.new,

  "#func": d.obj("Utilities for documenting Jsonnet methods (functions of objects)"),
  func:: {
    "#new": d.fn("new creates a new function, optionally with description and arguments", [d.arg("help", d.T.string), d.arg("args", d.T.array)]),
    new(help='', args=[]):: { 'function': {
      help: help,
      args: args,
    } },

    "#withHelp": d.fn("The `withHelp` modifier overrides the help text of that function", [d.arg("help", d.T.string)]),
    withHelp(help):: {'function'+: {
      help: help,
    }}
  },

  "#fn": self.func["#new"] + d.func.withHelp("`fn` is a shorthand for `func.new`"),
  fn:: self.func.new,

  "#argument": d.obj("Utilities for creating function arguments"),
  argument:: {
    "#new": d.fn("new creates a new function argument, taking the name, the type and optionally a default value", [d.arg("name", d.T.string), d.arg("type", d.T.string), d.arg("default", d.T.any)]),
    new(name, type, default=null): {
      name: name,
      type: type,
      default: default,
    },
  },
  "#arg": self.argument["#new"] + self.func.withHelp("`arg` is a shorthand for `argument.new`"),
  arg:: self.argument.new,

  // T contains constants for the Jsonnet types
  T:: {
    string: 'string',
    number: 'number',
    bool: 'bool',
    object: 'object',
    array: 'array',
    any: 'any',
    func: 'function',
  },
}

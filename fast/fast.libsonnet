local lib = {
  load(pkg):: {
    [key]:
      if std.startsWith(key, "#") then pkg[key]
      // else if std.objectHasAll(pkg[key], 'function') then
      //   pkg[key]
      // else if std.objectHasAll(pkg[key], 'object') then
      //   pkg[key] { object+: { fields: $.load(pkg[key].object.fields) } }
      else
        $.load(pkg[key])
    for key in std.objectFieldsAll(pkg)
    if std.isObject(pkg[key])
  },
};

local d = import 'doc-util/main.libsonnet';
local data = {
  '#': d.pkg('data', '', ''),

  '#tom': d.fn('yyay', [d.arg("name", d.T.string)]),
  tom(name): 'hi',

  '#foo':: d.obj('bar'),
  foo:: {
    '#baz': d.fn('foobar'),
    baz(): 'hi',
  },

  rock:: {
    '#stone': d.fn('hard', [d.arg("yay", d.T.string)]),
    stone(yay): 'hi',
  }
};

std.prune(lib.load(data))

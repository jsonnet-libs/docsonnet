local d = {
  new(name, url): {
    name: name,
    'import': url,
  },
  withAPI(obj): {
    api: obj.object,
  },
  withHelp(help): {
    help: help,
  },

  obj: {
    // obj creates an object, optionally with fields
    new(help='', fields={}): { object: {
      help: help,
      fields: fields,
    } },
  },
  Obj: self.obj.new,

  // fn creates a function, optionally with arguments
  fn: {
    new(help='', args=[]): { 'function': {
      help: help,
      args: args,
    } },
  },
  Fn: self.fn.new,

  // arg creates a function argument of given name,
  // type and optionally a default value
  arg: {
    new(name, type, default=null): {
      name: name,
      type: type,
      default: default,
    },
  },
  Arg: self.arg.new,

  // T contains constants for the Jsonnet types
  T: {
    string: 'string',
    number: 'number',
    bool: 'bool',
    object: 'object',
    array: 'array',
    any: 'any',
    func: 'function',
  },
};

local root = d.Obj('grafana.libsonnet is the offical Jsonnet library for Grafana', {
  new: d.Fn('new returns Grafana resources with sane defaults'),
  addConfig: d.Fn('addConfig adds config entries to grafana.ini', [
    d.Arg('config', d.T.object),
  ]),
  datasource: d.Obj('ds-util makes creating datasources easy', {
    new: d.Fn('new creates a new datasource', [
      d.Arg('name', d.T.string),
      d.Arg('type', d.T.string),
    ]),
    sheesh: d.Obj("sheesh is sheeshing around", {
      shit: d.Fn("enough sheesh", [
        d.Arg("lel", d.T.any),
        d.Arg("lol", d.T.func),
      ])
    })
  }),
});

d.new('grafana', 'github.com/sh0rez/grafana.libsonnet')
+ d.withAPI(std.prune(root))
+ d.withHelp("`grafana.libsonnet` is the offical Jsonnet package for using Grafana with Kubernetes")

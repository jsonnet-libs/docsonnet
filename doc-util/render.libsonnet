{
  local root = self,

  render(obj):
    assert std.isObject(obj) && '#' in obj : 'error: object is not a docsonnet package';
    local package = self.package(obj);
    package.toFiles(),

  findPackages(obj, path=[]): {
    local find(obj, path, parentWasPackage=true) =
      std.foldl(
        function(acc, k)
          acc
          + (
            // If matches a package, return it
            if '#' in obj[k]
            then [root.package(obj[k], path + [k], parentWasPackage)]
            // If matches a package but warn if also has an object docstring
            else if '#' in obj[k] && '#' + k in obj
            then std.trace(
              'warning: %s both defined as object and package' % k,
              [root.package(obj[k], path + [k], parentWasPackage)]
            )
            // If not, keep looking
            else find(obj[k], path + [k], parentWasPackage=false)
          ),
        std.filter(
          function(k)
            !std.startsWith(k, '#')
            && std.isObject(obj[k]),
          std.objectFieldsAll(obj)
        ),
        []
      ),

    packages: find(obj, path),

    hasPackages(): std.length(self.packages) > 0,

    toIndex():
      if self.hasPackages()
      then
        std.join('\n', [
          '* ' + p.link
          for p in self.packages
        ])
        + '\n'
      else '',

    toFiles():
      std.foldl(
        function(acc, p)
          acc
          + { [p.path]: p.toString() }
          + p.packages.toFiles(),
        self.packages,
        {}
      ),
  },

  package(obj, path=[], parentWasPackage=true): {
    local this = self,
    local doc = obj['#'],

    packages: root.findPackages(obj, path),
    fields: root.fields(obj),

    path:
      std.join(
        '/',
        path
      )
      + (if self.packages.hasPackages()
         then '/index.md'
         else '.md'),

    link:
      '[%s](%s)' % [
        std.join(
          '.',
          (if parentWasPackage  // if parent object was a package
           then path[1:]  // then strip first item from path
           else path)
        ),
        self.path,
      ],

    toFiles():
      { 'README.md': this.toString() }
      + self.packages.toFiles(),

    toString():
      std.join(
        '\n',
        ['# ' + doc.name + '\n']
        + (if std.get(doc, 'help', '') != ''
           then [doc.help]
           else [])
        + (if self.packages.hasPackages()
           then [
             '## Subpackages\n\n'
             + self.packages.toIndex(),
           ]
           else [])
        + (if self.fields.hasFields()
           then [
             '## Index\n\n'
             + self.fields.toIndex()
             + '\n## Fields\n'
             + self.fields.toString(),
           ]
           else [])
      ),
  },

  fields(obj, path=[]): {
    values: root.findValues(obj, path),
    functions: root.findFunctions(obj, path),
    objects: root.findObjects(obj, path),

    hasFields():
      std.any([
        self.values.hasFields(),
        self.functions.hasFields(),
        self.objects.hasFields(),
      ]),

    toIndex():
      std.join('', [
        self.functions.toIndex(),
        self.objects.toIndex(),
      ]),

    toString():
      std.join('', [
        self.values.toString(),
        self.functions.toString(),
        self.objects.toString(),
      ]),
  },

  findObjects(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('object', obj),
        std.objectFieldsAll(obj)
      ),

    local undocumentedKeys =
      std.filter(
        function(k)
          std.all([
            !std.startsWith(k, '#'),
            std.isObject(obj[k]),
            !('#' + k in obj),  // not documented in parent
            !('#' in obj[k]),  // not a sub package
          ]),
        std.objectFieldsAll(obj)
      ),

    objects:
      std.foldl(
        function(acc, k)
          acc + [
            root.obj(
              root.util.realkey(k),
              obj[k],
              obj[root.util.realkey(k)],
              path,
            ),
          ],
        keys,
        []
      )
      + std.foldl(
        function(acc, k)
          local o = root.obj(
            k,
            { object: { help: '' } },
            obj[k],
            path,
          );
          acc
          + (if o.fields.hasFields()
             then [o]
             else []),
        undocumentedKeys,
        []
      ),

    hasFields(): std.length(self.objects) > 0,

    toIndex():
      if self.hasFields()
      then
        std.join('', [
          std.join(
            '',
            [' ' for d in std.range(0, (std.length(path) * 2) - 1)]
            + ['* ', f.link]
            + ['\n']
            + (if f.fields.hasFields()
               then [f.fields.toIndex()]
               else [])
          )
          for f in self.objects
        ])
      else '',

    toString():
      if self.hasFields()
      then
        std.join('', [
          o.toString()
          for o in self.objects
        ])
      else '',
  },

  obj(name, doc, obj, path): {
    fields: root.fields(obj, path + [name]),

    path: std.join('.', path + [name]),
    fragment: root.util.fragment(std.join('', path + [name])),
    link: '[`obj %s`](#obj-%s)' % [name, self.fragment],

    toString():
      std.join(
        '\n',
        [root.util.title('obj ' + self.path, std.length(path) + 2)]
        + (if doc.object.help != ''
           then [doc.object.help]
           else [])
        + [self.fields.toString()]
      ),
  },

  findFunctions(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('function', obj),
        std.objectFieldsAll(obj)
      ),

    functions:
      std.foldl(
        function(acc, k)
          acc + [
            root.func(
              root.util.realkey(k),
              obj[k],
              path,
            ),
          ],
        keys,
        []
      ),

    hasFields(): std.length(self.functions) > 0,

    toIndex():
      if self.hasFields()
      then
        std.join('', [
          std.join(
            '',
            [' ' for d in std.range(0, (std.length(path) * 2) - 1)]
            + ['* ', f.link]
            + ['\n']
          )
          for f in self.functions
        ])
      else '',

    toString():
      if self.hasFields()
      then
        std.join('', [
          f.toString()
          for f in self.functions
        ])
      else '',
  },

  func(name, doc, path): {
    path: std.join('.', path + [name]),
    fragment: root.util.fragment(std.join('', path + [name])),
    link: '[`fn %s(%s)`](#fn-%s)' % [name, self.args, self.fragment],

    args: std.join(', ', [
      if arg.default != null
      then std.join('=', [
        arg.name,
        std.manifestJsonEx(arg.default, '', ''),
      ])
      else arg.name
      for arg in doc['function'].args
    ]),

    enums: std.join('', [
      if arg.enums != null
      then '\n\nAccepted values for `%s` are ' % arg.name
           + std.join(', ', [
             std.manifestJsonEx(item, '', '')
             for item in arg.enums
           ])
      else ''
      for arg in doc['function'].args
    ]),

    toString():
      std.join('\n', [
        root.util.title('fn ' + self.path, std.length(path) + 2),
        |||
          ```jsonnet
          %(name)s(%(args)s)
          ```
        ||| % [name, self.args],
        doc['function'].help,
        self.enums,
      ]),
  },

  findValues(obj, path=[]): {
    local keys =
      std.filter(
        root.util.filter('value', obj),
        std.objectFieldsAll(obj)
      ),

    values:
      std.foldl(
        function(acc, k)
          acc + [
            root.val(
              root.util.realkey(k),
              obj[k],
              obj[root.util.realkey(k)],
              path,
            ),
          ],
        keys,
        []
      ),

    hasFields(): std.length(self.values) > 0,

    toString():
      if self.hasFields()
      then
        std.join('\n', [
          '* ' + f.toString()
          for f in self.values
        ]) + '\n'
      else '',
  },

  val(name, doc, obj, path): {
    toString():
      std.join(' ', [
        '`%s`' % std.join('.', path + [name]),
        '(`%s`):' % doc.value.type,
        '`"%s"`' % obj,
        '-',
        doc.value.help,
      ]),
  },

  util: {
    realkey(key):
      assert std.startsWith(key, '#') : 'Key %s not a docstring key' % key;
      key[1:],
    title(title, depth=0):
      std.join(
        '',
        ['\n']
        + ['#' for i in std.range(0, depth)]
        + [' ', title, '\n']
      ),
    fragment(title):
      std.asciiLower(
        std.strReplace(
          std.strReplace(title, '.', '')
          , ' ', '-'
        )
      ),
    filter(type, obj):
      function(k)
        std.all([
          std.startsWith(k, '#'),
          std.isObject(obj[k]),
          type in obj[k],
          root.util.realkey(k) in obj,
        ]),
  },
}

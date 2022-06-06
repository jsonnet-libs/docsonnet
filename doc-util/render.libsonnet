{
  local root = self,

  package(p):: {
    name: p.name,
    doc: |||
      # package %(name)s

      ```jsonnet
      local %(name)s = import '%(import)s/main.libsonnet';
      ```

      %(help)s

    ||| % p,
  },

  headerPrefix(depth):: std.join('', [
    '#'
    for d in std.range(0, depth + 2)
  ]),

  func(k, v, depth):: {
    local data = {
      name: k[1:],
      nameLower: std.asciiLower(self.name),
      help: v['function'].help,
      args: std.join(', ', [
        if arg.default != null
        then arg.name + '=' + arg.default
        else arg.name
        for arg in v['function'].args
      ]),
      headerPrefix: root.headerPrefix(depth),
    },
    index: '[`fn %(name)s(%(args)s)`](#fn-%(nameLower)s)' % data,
    doc: |||
      %(headerPrefix)s fn %(name)s

      ```ts
      %(name)s(%(args)s)
      ```

      %(help)s

    ||| % data,
  },

  obj(k, v, obj, depth):: {
    local spaces = std.join('', [
      ' '
      for d in std.range(0, (depth * 2) + 1)
    ]),
    local processed = root.process(obj, depth + 1),
    local data = {
      name: k[1:],
      nameLower: std.asciiLower(self.name),
      help: v.object.help,
      subs: root.renderFields(processed, false),
      index: root.index(processed, spaces),
      headerPrefix: root.headerPrefix(depth),
    },
    index: '[`obj %(name)s`](#obj-%(nameLower)s)\n%(index)s' % data,
    doc: |||
      %(headerPrefix)s obj %(name)s

      %(help)s
      %(subs)s

    ||| % data,
  },

  process(obj, depth=0)::
    std.foldl(
      function(acc, k)
        acc +
        if std.startsWith(k, '#')
        then (
          local realKey = k[1:];
          if 'function' in obj[k]
          then { fields+: [root.func(k, obj[k], depth)] }
          else if 'object' in obj[k]
          then { objects+: [root.obj(k, obj[k], obj[realKey], depth)] }
          else { package: root.package(obj[k]) }
        )
        else (
          if '#' + k in obj
          then {}  // if has docs, do not search for package
          else
            if std.isObject(obj[k])
            then {
              package+: {
                subpackages+: [root.process(obj[k])],
              },
            }
            else {}
        ),
      std.objectFieldsAll(obj),
      {}
    ),

  index(p, spaces='')::
    (
      if 'fields' in p
      then std.join('\n', [
        spaces + '* ' + f.index
        for f in p.fields
      ]) + '\n'
      else ''
    )
    + (
      if 'objects' in p
      then std.join('\n', [
        spaces + '* ' + f.index
        for f in p.objects
      ])
      else ''
    ),

  renderFields(p, headers=true)::
    (
      if 'fields' in p
      then
        (if headers then '## Fields\n\n' else '')
        + std.join('\n', [
          f.doc
          for f in p.fields
        ])
      else ''
    )
    + (
      if 'objects' in p
      then std.join('\n', [
        f.doc
        for f in p.objects
      ])
      else ''
    ),

  renderPackage(p, headers=true, index=true)::
    (if 'package' in p
     then p.package.doc
          + (
            if 'subpackages' in p.package
            then
              (if headers then '## Subpackages\n\n' else '')
              + std.join('\n', [
                '* [%(name)s](%(path)s)' % {
                  name: sub.package.name,
                  path: p.package.name + '/' + sub.package.name + '.md',
                }
                for sub in p.package.subpackages
                if 'package' in sub
              ]) + '\n\n'
            else ''
          )
     else '')
    + (if headers then '## Index\n\n' else '')
    + (if index then root.index(p) else '')
    + root.renderFields(p, headers),

  render(p, prefix='')::
    if 'package' in p
    then
      local subs =
        if 'subpackages' in p.package
        then [
          root.render(sub, prefix=prefix + p.package.name + '/')
          for sub in p.package.subpackages
        ]
        else [];
      {
        [prefix + p.package.name + '.md']: root.renderPackage(p),
      }
      + std.foldl(
        function(acc, sub)
          acc + sub,
        subs,
        {}
      )
    else {},
}

{
  local root = self,

  templates: {
    package: |||
      # package %(name)s

      ```jsonnet
      local %(name)s = import '%(import)s/%(filename)s';
      ```

      %(help)s
    |||,

    index: |||
      ## Index

      %s
    |||,

    sectionTitle: '%(abbr)s %(prefix)s%(name)s',

    sectionLink: '* [`%(abbr)s %(linkName)s`](#%(link)s)',

    section: |||
      %(headerDepth)s %(title)s

      %(content)s
    |||,
  },

  renderSectionTitle(section, prefixes)::
    root.templates.sectionTitle % {
      name: section.name,
      abbr: section.type.abbr,
      prefix: std.join('.', prefixes) + if std.length(prefixes) > 0 then '.' else '',
    },

  renderSection(section, depth=0, prefixes=[])::
    root.templates.section % {
      headerDepth: std.join('', [
        '#'
        for d in std.range(0, depth + 2)
      ]),
      title: root.renderSectionTitle(
        section,
        prefixes,
      ),
      content: section.content,
    }
    + (
      if std.length(section.subSections) > 0
      then
        std.join('\n', [
          root.renderSection(subsection, depth + 1, prefixes + [section.name])
          for subsection in section.subSections
        ])
      else ''
    ),

  renderPackage(package)::
    (root.templates.package % package)
    + (
      if std.length(package.subPackages) > 0
      then
        '## Subpackages\n\n'
        + std.join('\n', [
          '* [%(name)s](%(path)s)' % {
            name: sub.name,
            path: package.name + '/' + sub.name + '.md',
          }
          for sub in package.subPackages
        ]) + '\n\n'
      else ''
    )
    + (root.templates.index % root.index(package.sections))
    + '\n## Fields\n\n'
    + std.join('\n', [
      root.renderSection(section)
      for section in package.sections
    ]),

  index(sections, depth=0, prefixes=[])::
    std.join('\n', [
      std.join('', [
        ' '
        for d in std.range(0, (depth * 2) - 1)
      ])
      + (root.templates.sectionLink % {
           abbr: section.type.abbr,
           linkName: section.linkName,
           link:
             std.asciiLower(
               std.strReplace(
                 std.strReplace(root.renderSectionTitle(section, prefixes), '.', '')
                 , ' ', '-'
               )
             ),
         })
      + (
        if std.length(section.subSections) > 0
        then '\n' + root.index(section.subSections, depth + 1, prefixes + [section.name])
        else ''
      )
      for section in sections
    ]),

  section(key, doc, obj, depth):: {
    name: std.strReplace(key, '#', ''),

    subSections:
      if std.isObject(obj)
      then root.process(obj, depth=depth + 1).sections
      else [],

    type:
      if 'function' in doc
      then { full: 'function', abbr: 'fn' }
      else if 'object' in doc
      then { full: 'object', abbr: 'obj' }
      else if std.isObject(obj)
      then { full: 'object', abbr: 'obj' }
      else { full: '', abbr: '' },

    abbr: self.type.abbr,

    doc:
      if self.type.full in doc
      then doc[self.type.full]
      else { help: '' },

    help: self.doc.help,

    args:
      if 'args' in self.doc
      then std.join(', ', [
        if arg.default != null
        then arg.name + '=' + arg.default
        else arg.name
        for arg in self.doc.args
      ])
      else '',

    linkName:
      if 'args' in self.doc
      then '%(name)s(%(args)s)' % self
      else self.name,

    contentTemplate:
      if self.type.full == 'function'
      then '```ts\n%(name)s(%(args)s)\n```\n\n%(help)s'
      else |||
        %(help)s
      |||,

    content: self.contentTemplate % self,
  },

  process(obj, filename='', depth=0)::
    std.foldl(
      function(acc, key)
        acc +
        // Package definition
        if key == '#'
        then obj[key] { filename: filename }

        // Field definition
        else if std.startsWith(key, '#')
        then (
          local realKey = key[1:];
          if 'function' in obj[key]
          then {
            functionSections+: [root.section(key, obj[key], obj[realKey], depth)],
          }
          else if 'object' in obj[key]
          then {
            objectSections+: [root.section(key, obj[key], obj[realKey], depth)],
          }
          else {}
        )

        // subPackage definition
        else if std.isObject(obj[key]) && '#' in obj[key]
        then {
          subPackages+: [root.process(obj[key])],
        }

        // undocumented object
        else if std.isObject(obj[key]) && !('#' + key in obj)
        then (
          local section = root.section(key, {}, obj[key], depth);
          // only add if has documented subSections
          if std.length(section.subSections) > 0
          then { objectSections+: [section] }
          else {}
        )

        else {},
      std.objectFieldsAll(obj),
      {
        sections: self.functionSections + self.objectSections,
        functionSections: [],
        objectSections: [],
        subPackages: [],
      }
    ),

  renderFiles(package, prefix='')::
    {
      [if prefix == '' then 'README.md' else prefix + package.name + '.md']: root.renderPackage(package),
    }
    + std.foldl(
      function(acc, sub)
        acc + sub,
      [
        root.renderFiles(sub, prefix=prefix + package.name + '/')
        for sub in package.subPackages
      ],
      {}
    ),

  render(obj, filename):
    self.renderFiles(self.process(obj, filename)),
}

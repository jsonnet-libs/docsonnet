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

    value: '* `%(prefix)s%(name)s` (`%(type)s`): `"%(value)s"` - %(help)s',

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

  renderValues(values, prefixes=[])::
    if std.length(values) > 0
    then
      std.join('\n', [
        root.templates.value
        % value {
          prefix: std.join('.', prefixes) + if std.length(prefixes) > 0 then '.' else '',
        }
        for value in values
      ]) + '\n'
    else '',

  renderSections(sections, depth=0, prefixes=[])::
    if std.length(sections) > 0
    then
      std.join('\n', [
        root.templates.section
        % {
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
        + root.renderValues(
          section.values,
          prefixes + [section.name]
        )
        + root.renderSections(
          section.subSections,
          depth + 1,
          prefixes + [section.name]
        )
        for section in sections
      ])
    else '',

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
    + root.renderValues(package.values)
    + root.renderSections(package.sections),

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

    local processed =
      if std.isObject(obj)
      then root.process(obj, depth=depth + 1)
      else { sections: [], values: [] },

    subSections: processed.sections,

    values: processed.values,

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
      else if self.help != ''
      then '%(help)s\n'
      else '',

    content: self.contentTemplate % self,
  },

  value(key, doc, obj):: {
    name: std.strReplace(key, '#', ''),
    type: doc.value.type,
    help: doc.value.help,
    value: obj,
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
          if 'value' in std.trace(key, obj[key])
          then {
            values+: [root.value(key, obj[key], obj[realKey])],
          }
          else if 'function' in obj[key]
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
          // only add if has documented subSections or values
          if std.length(section.subSections) > 0
             || std.length(section.values) > 0
          then { objectSections+: [section] }
          else {}
        )

        else {},
      std.objectFieldsAll(obj),
      {
        functionSections: [],
        objectSections: [],

        sections:
          self.functionSections
          + self.objectSections,
        subPackages: [],
        values: [],
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

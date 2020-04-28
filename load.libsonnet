{
  // reshape converts the Jsonnet structure to the one used by docsonnet:
  // - put fields into an `api` key
  // - put subpackages into `sub` key
  reshape(pkg)::
    local aux(old, key) =
      if key == '#' then
        old
      else if std.objectHas(pkg[key], '#') then
        old { sub+: { [key]: $.package(pkg[key]) } }
      else
        old { api+: { [key]: pkg[key] } };

    std.foldl(aux, std.objectFields(pkg), {})
    + pkg['#'],

  // fillObjects creates docsonnet objects from Jsonnet ones,
  // also filling those that have been specified explicitely
  fillObjects(api)::
    local aux(old, key) =
      if std.startsWith(key, '#') then
        old { [key]: api[key] }
      else if std.isObject(api[key]) && std.length(std.objectFields(api[key])) > 0 then
        old { ['#' + key]+: { object+: {
          fields: api[key],
        } } }
      else old;

    std.foldl(aux, std.objectFields(api), {}),

  // clean removes all hashes from field names
  clean(api):: {
    [std.lstripChars(key, '#')]:
      if std.isObject(api[key]) then $.clean(api[key])
      else api[key]
    for key in std.objectFields(api)
  },

  cleanNonObj(api):: {
    [key]:
      if std.startsWith(key, "#") then api[key]
      else if std.isObject(api[key]) then $.cleanNonObj(api[key])
      else api[key]
    for key in std.objectFieldsAll(api)
    if std.isObject(api[key])
  },

  // package loads docsonnet from a Jsonnet package
  package(pkg)::
    local cleaned = self.cleanNonObj(pkg);
    local reshaped = self.reshape(cleaned);
    local filled =
      if std.objectHas(reshaped, 'api')
      then reshaped { api: $.fillObjects(reshaped.api) }
      else reshaped;
    self.clean(filled),
    // reshaped,
}

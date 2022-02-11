local lib = {
  scan(obj)::
    local aux(old, key) =
      if std.startsWith(key, '#') then
        true
      else if std.isObject(obj[key]) then
        old || $.scan(obj[key])
      else old;
    std.foldl(aux, std.objectFieldsAll(obj), false),

  load(pkg)::
    local aux(old, key) =
      if !std.isObject(pkg[key]) then
        old
      else if std.objectHasAll(pkg, '#' + key) && pkg['#' + key] == 'ignore' then
        old
      else if std.startsWith(key, '#') then
        old { [key]: pkg[key] }
      else if self.scan(pkg[key]) then
        old { [key]: $.load(pkg[key]) }
      else old;

    std.foldl(aux, std.objectFieldsAll(pkg), {}),
};


lib.load(std.extVar('main'))

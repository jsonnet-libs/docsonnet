local d = import 'doc-util/main.libsonnet';

local pkg(name, url) = {
  NAME:: {
    TYPE:: { help: "" },
    name: name,
    'import': url,
    help: self.TYPE.help,
  },
  "#": self.NAME,
};

local dType(kind) = function(name) {
  // NAME is used to mix into the docsonnet field without knowing it's name
  NAME:: {
    // TYPE is used to mix into {function,object,value} without knowing what it is
    TYPE:: {},
    [kind]: self.TYPE,
  },
  ['#' + name]: self.NAME,
};

local fn = dType('function'),
      obj = dType('object');

local val(name, type, default=null) = dType('value') + {
  NAME+: { TYPE+: {
    type: type,
    default: default,
  } },
};

local arg(name, type, default=null) = {
  NAME+: { TYPE+: {
    args+: [{
      name: name,
      type: type,
      default: default,
    }],
  } },
};

local desc(help) = {
  NAME+: { TYPE+: {
    help: help,
  } },
};

local string = 'string',
      object = 'object';

{}

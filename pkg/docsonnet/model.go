package docsonnet

// Package represents a Jsonnet package, having an API (list of Fields) and
// perhaps subpackages
type Package struct {
	Name   string `json:"name"`
	Import string `json:"import"`
	Help   string `json:"help"`

	API Fields             `json:"api,omitempty"`
	Sub map[string]Package `json:"sub,omitempty"`
}

// Object represents a Jsonnet object, a list of key-value fields
type Object struct {
	Name string `json:"-"`
	Help string `json:"help"`

	// children
	Fields Fields `json:"fields"`
}

// Function represents a Jsonnet function, a named construct that takes
// arguments
type Function struct {
	Name string `json:"-"`
	Help string `json:"help"`

	Args []Argument `json:"args,omitempty"`
}

// Argument is a function argument, optionally also having a default value
type Argument struct {
	Type    Type        `json:"type"`
	Name    string      `json:"name"`
	Default interface{} `json:"default"`
}

// Value is a value of any other type than the special Object and Function types
type Value struct {
	Name string `json:"-"`
	Help string `json:"help"`

	Type    Type        `json:"type"`
	Default interface{} `json:"default"`
}

// Type is a Jsonnet type
type Type string

const (
	TypeString = "string"
	TypeNumber = "number"
	TypeBool   = "boolean"
	TypeObject = "object"
	TypeArray  = "array"
	TypeAny    = "any"
	TypeFunc   = "function"
)

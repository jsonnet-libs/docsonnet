package docsonnet

import (
	"encoding/json"
	"errors"
)

// Field represents any field of an object.
type Field struct {
	// Function value
	Function *Function `json:"function,omitempty"`
	// Object value
	Object *Object `json:"object,omitempty"`
	// Any other value
	Value *Value `json:"value,omitempty"`
}

func (o *Field) UnmarshalJSON(data []byte) error {
	type fake Field

	var f fake
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}

	switch {
	case f.Function != nil:
		o.Function = f.Function
	case f.Object != nil:
		o.Object = f.Object
	case f.Value != nil:
		o.Value = f.Value
	default:
		return errors.New("field has no value")
	}

	return nil
}

func (o Field) MarshalJSON() ([]byte, error) {
	if o.Function == nil && o.Object == nil && o.Value == nil {
		return nil, errors.New("field has no value")
	}

	type fake Field
	return json.Marshal(fake(o))
}

// Fields is a list of fields
type Fields map[string]Field

func (fPtr *Fields) UnmarshalJSON(data []byte) error {
	if *fPtr == nil {
		*fPtr = make(Fields)
	}
	f := *fPtr

	tmp := make(map[string]Field)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for k, v := range tmp {
		switch {
		case v.Function != nil:
			v.Function.Name = k
		case v.Object != nil:
			v.Object.Name = k
		case v.Value != nil:
			v.Value.Name = k
		}
		f[k] = v
	}

	return nil
}

package swagger

import (
	"encoding/json"
	"reflect"
)

// Ref represents a json reference
type Ref string

// MarshalJSON marshal this to JSON
func (r Ref) MarshalJSON() ([]byte, error) {
	if r == "" {
		return []byte("{}"), nil
	}
	v := map[string]interface{}{"$ref": string(r)}
	return json.Marshal(v)
}

// UnmarshalJSON unmarshal this from JSON
func (r *Ref) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	if vv, ok := v["$ref"]; ok {
		if str, ok := vv.(string); ok {
			*r = Ref(str)
		}
	}
	return nil
}

// Definitions contains the models explicitly defined in this spec
// An object to hold data types that can be consumed and produced by operations.
// These data types can be primitives, arrays or models.
//
// For more information: http://goo.gl/8us55a#definitionsObject
type Definitions map[string]Schema

// SecurityDefinitions a declaration of the security schemes available to be used in the specification.
// This does not enforce the security schemes on the operations and only serves to provide
// the relevant details for each scheme.
//
// For more information: http://goo.gl/8us55a#securityDefinitionsObject
type SecurityDefinitions map[string]*SecurityScheme

// StringOrArray represents a value that can either be a string
// or an array of strings. Mainly here for serialization purposes
type StringOrArray struct {
	Single string
	Multi  []string
}

// UnmarshalJSON unmarshals this string or array object from a JSON array or JSON string
func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	if data[0] == '[' {
		var parsed []string
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		s.Multi = parsed
		return nil
	}

	var parsed string
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	s.Single = parsed
	return nil
}

// MarshalJSON converts this string or array to a JSON array or JSON string
func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if s.Single != "" {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

// SchemaOrArray represents a value that can either be a Schema
// or an array of Schema. Mainly here for serialization purposes
type SchemaOrArray struct {
	Single *Schema
	Multi  []Schema
}

// MarshalJSON converts this schema object or array into JSON structure
func (s SchemaOrArray) MarshalJSON() ([]byte, error) {
	// fmt.Println("marshalli")
	if s.Single != nil {
		return json.Marshal(s.Single)
	}
	return json.Marshal(s.Multi)
}

// UnmarshalJSON converts this schema object or array from a JSON structure
func (s *SchemaOrArray) UnmarshalJSON(data []byte) error {
	if len(data) < 3 {
		return nil
	}
	if data[0] == '[' {
		var parsed []Schema
		if err := json.Unmarshal(data, &parsed); err != nil {
			return err
		}
		s.Multi = parsed
		return nil
	}

	var parsed Schema
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	if reflect.DeepEqual(Schema{}, parsed) {
		return nil
	}
	s.Single = &parsed
	return nil
}

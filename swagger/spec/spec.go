package spec

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/assets"
	"github.com/casualjim/go-swagger/swagger/jsonschema"
)

// MustLoadSwagger20Schema panics when Swagger20Schema returns an error
func MustLoadSwagger20Schema() *jsonschema.Document {
	d, e := Swagger20Schema()
	if e != nil {
		panic(e)
	}
	return d
}

// Swagger20Schema loads the swagger 2.0 schema from the embedded asses
func Swagger20Schema() (*jsonschema.Document, error) {

	b, err := assets.Asset("2.0/schema.json")
	if err != nil {
		return nil, err
	}
	loader := jsonschema.NewLoader(bytes.NewBuffer(b), "http://swagger.io/v2/schema.json#")

	return jsonschema.Load(loader)
}

// Document represents a swagger spec document
type Document struct {
	specAnalyzer
	specSchema *jsonschema.Document
	data       map[string]interface{}
	spec       *swagger.Spec
}

// New creates a new shema document
func New(data json.RawMessage, version string) (*Document, error) {
	if version == "" {
		version = "2.0"
	}
	if version != "2.0" {
		return nil, fmt.Errorf("spec version %q is not supported", version)
	}

	specSchema := MustLoadSwagger20Schema()

	spec := new(swagger.Spec)
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}

	var v map[string]interface{}
	json.Unmarshal(data, &v)

	d := &Document{
		specAnalyzer: specAnalyzer{
			spec:        spec,
			consumes:    make(map[string]struct{}),
			produces:    make(map[string]struct{}),
			authSchemes: make(map[string]struct{}),
			operations:  make(map[string]map[string]*swagger.Operation),
		},
		specSchema: specSchema,
		data:       v,
		spec:       spec,
	}
	d.initialize()
	return d, nil
}

// BasePath the base path for this spec
func (d *Document) BasePath() string {
	return d.spec.BasePath
}

// Version returns the version of this spec
func (d *Document) Version() string {
	return d.spec.Swagger
}

// Validate validates this spec document
func (d *Document) Validate() *jsonschema.ValidationResult {
	return d.specSchema.Validate(d.data)
}

package jsonschema

import "github.com/iancoleman/orderedmap"

func NewSchema(types ...string) *Schema {
	typeName := "string"
	if len(types) > 0 {
		typeName = types[0]
	}

	var schema = new(Schema)
	schema.Type = typeName

	if typeName == "object" {
		schema.Properties = orderedmap.New()
	}
	return schema
}

func NewSchemaSetItems(typeName string) *Schema {
	var schema = NewSchema("array")
	schema.Items = NewSchema(typeName)
	return schema
}

func (t *Schema) IsObj() bool {
	return t.Type == "object"
}

func (t *Schema) IsArray() bool {
	return t.Type == "array"
}

func (t *Schema) IsNull() bool {
	return t.Type == "null"
}

func (t *Schema) IsSpread() bool {
	return !t.IsNormal()
}

func (t *Schema) IsNormal() bool {
	return !t.IsObj() && !t.IsArray() && !t.IsNull()
}

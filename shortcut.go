package jsonschema

func NewSchema(types ...string) *Schema {
	typeName := "string"
	if len(types) > 0 {
		typeName = types[0]
	}

	var schema = new(Schema)
	schema.Type = typeName

	if typeName == "object" {
		schema.Properties = NewProperties()
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

func (t *Schema) AddMeta(key string, value interface{}) {
	if t.MetaData == nil {
		t.MetaData = make(map[string]interface{})
	}
	t.MetaData[key] = value
}

func (t *Schema) GetMeta(key string) (interface{}, bool) {
	if t.MetaData == nil {
		return nil, false
	}
	v, ok := t.MetaData[key]
	return v, ok
}

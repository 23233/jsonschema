package jsonschema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetSchemaMapByPointer(t *testing.T) {
	schemaJSON := `
        {
            "type": "object",
            "properties": {
                "foo": {
                    "type": "string"
                },
                "bar": {
                    "type": "array",
                    "items": [
                        {
                            "type": "number"
                        },
                        {
                            "type": "object",
                            "properties": {
                                "baz": {
                                    "type": "string"
                                }
                            },
                            "required": ["baz"]
                        }
                    ]
                }
            },
            "required": ["foo", "bar"]
        }
    `
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(schemaJSON), &schema); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	tests := []struct {
		pointer string
		expect  interface{}
		hasErr  bool
	}{
		// valid cases
		{"", nil, true},
		{"/", schema, false},
		{"/foo", map[string]interface{}{"type": "string"}, false},
		{"/bar", map[string]interface{}{"type": "array", "items": []interface{}{
			map[string]interface{}{"type": "number"},
			map[string]interface{}{"type": "object", "properties": map[string]interface{}{
				"baz": map[string]interface{}{"type": "string"},
			}, "required": []interface{}{"baz"}},
		}}, false},
		{"/bar/0", map[string]interface{}{"type": "number"}, false},
		{"/bar/1", map[string]interface{}{"type": "object", "properties": map[string]interface{}{
			"baz": map[string]interface{}{"type": "string"},
		}, "required": []interface{}{"baz"}}, false},
		{"/bar/1/baz", map[string]interface{}{"type": "string"}, false},
		// 暂时不支持 - ~ 这种操作符
		{"/bar/-", nil, true},
		{"/bar/-/baz", nil, true},
		{"/bar/-/foo", nil, true},

		// invalid cases
		{"/foo/bar", nil, true},
		{"/bar/foo", nil, true},
		{"/baz", nil, true},
		{"/bar/2", nil, true},
		{"/bar/1/foo", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.pointer, func(t *testing.T) {
			actual, err := GetSchemaMapByPointer(schema, tt.pointer)
			if tt.hasErr {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(actual, tt.expect) {
					t.Errorf("got: %#v, want: %#v", actual, tt.expect)
				}
			}
		})
	}

	// 测试$ref
	refSchema := `{"$defs":{"ModelIndex":{"additionalProperties":false,"properties":{"field_name":{"items":{"type":"string"},"type":"array"},"type":{"type":"string"}},"type":"object"},"RawSchema":{"type":"object","widget":"RawJsonTree"}},"$id":"https://resok.cn/s/schemas/model","$schema":"https://json-schema.org/draft/2020-12/schema","additionalProperties":false,"properties":{"backend":{"default":"mongodb","enum":["mongodb"],"type":"string"},"desc":{"type":"string"},"fieldsDefine":{"$ref":"#/$defs/RawSchema"},"group":{"type":"string"},"indexes":{"items":{"$ref":"#/$defs/ModelIndex"},"type":"array"},"title":{"type":"string"},"user_id":{"type":"string"}},"required":["fieldsDefine","title"],"title":"模型","type":"object"}`
	var refSchemaJSON map[string]interface{}
	if err := json.Unmarshal([]byte(refSchema), &refSchemaJSON); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	refTest := []struct {
		pointer string
		expect  interface{}
		hasErr  bool
	}{

		{"#/indexes", map[string]interface{}{"items": map[string]interface{}{"$ref": "#/$defs/ModelIndex"}, "type": "array"}, false},
		{"#/fieldsDefine", map[string]any{"type": "object", "widget": "RawJsonTree"}, false},
		// 为什么非要有items关键词 因为数组没有key
		{"#/indexes/items", map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"field_name": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "type": map[string]any{"type": "string"}}}, false},
		// 这里为什么不是 #/indexes/items/properties/type 呢? 因为#/indexes/items获取到的schema 已经是 ModelIndex了 下一个/type 就会自动判断object 去properties中找
		{"#/indexes/items/type", map[string]any{"type": "string"}, false},
	}
	for _, tt := range refTest {
		t.Run(tt.pointer, func(t *testing.T) {
			actual, err := GetSchemaMapByPointer(refSchemaJSON, tt.pointer)
			if tt.hasErr {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(actual, tt.expect) {
					t.Errorf("got: %#v, want: %#v", actual, tt.expect)
				}
			}
		})
	}

}

func TestSchemaHelper_GenAccessKeys(t *testing.T) {
	refSchema := `{"$defs":{"ModelIndex":{"additionalProperties":false,"properties":{"field_name":{"items":{"type":"string"},"type":"array"},"type":{"type":"string"}},"type":"object"},"RawSchema":{"type":"object","widget":"RawJsonTree"}},"$id":"https://resok.cn/s/schemas/model","$schema":"https://json-schema.org/draft/2020-12/schema","additionalProperties":false,"properties":{"backend":{"default":"mongodb","enum":["mongodb"],"type":"string"},"desc":{"type":"string"},"fieldsDefine":{"$ref":"#/$defs/RawSchema"},"group":{"type":"string"},"indexes":{"items":{"$ref":"#/$defs/ModelIndex"},"type":"array"},"title":{"type":"string"},"user_id":{"type":"string"}},"required":["fieldsDefine","title"],"title":"模型","type":"object"}`
	var refSchemaJSON map[string]interface{}
	if err := json.Unmarshal([]byte(refSchema), &refSchemaJSON); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}
	helper := NewSchemaHelper(refSchemaJSON)
	accessKeys := helper.GenAccessKeys()

	assert.Equal(t, len(accessKeys), len(refSchemaJSON["properties"].(map[string]any))+1)

	expected := []string{"desc", "fieldsDefine", "group", "indexes.*.field_name", "indexes.*.type", "title", "user_id", "backend"}
	for _, s := range expected {
		has := false
		for _, key := range accessKeys {
			if s == key {
				has = true
				break
			}
		}
		if !has {
			t.Errorf("accessKey 应该包含%s 但是未获取到 ", s)
			return
		}
	}
}

func TestFindDataByAccessKey(t *testing.T) {
	// 测试针对对象的情况
	obj := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}
	var expected any
	expected = "John"
	result := FindDataByAccessKey(obj, "name")
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	// 测试针对数组的情况
	arr := []interface{}{
		map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		map[string]interface{}{
			"name": "Mary",
			"age":  25,
		},
	}
	expected = "Mary"
	result = FindDataByAccessKey(arr, "1.name")
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	// 测试针对数组和对象的混合情况
	mixedData := map[string]interface{}{
		"name": "John",
		"age":  30,
		"pets": []interface{}{
			map[string]interface{}{
				"name": "Fluffy",
				"type": "cat",
			},
			map[string]interface{}{
				"name": "Fido",
				"type": "dog",
			},
		},
	}
	expected = "Fido"
	result = FindDataByAccessKey(mixedData, "pets.1.name")
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	// 测试找不到数据的情况
	expectedArr := make([]any, 0)
	resultArr := FindDataByAccessKey(mixedData, "pets.2.type")
	if !reflect.DeepEqual(resultArr, expectedArr) {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	// 测试带有通配符的情况
	mixedData = map[string]interface{}{
		"name": "John",
		"age":  30,
		"pets": []interface{}{
			map[string]interface{}{
				"name": "Fluffy",
				"type": "cat",
			},
			map[string]interface{}{
				"name": "Fido",
				"type": "dog",
			},
		},
		"favoriteFoods": map[string]interface{}{
			"breakfast": "pancakes",
			"lunch":     "sandwich",
			"dinner":    "pizza",
		},
	}
	expectedArr = []interface{}{"cat", "dog"}
	resultArr = FindDataByAccessKey(mixedData, "pets.*.type")
	if !reflect.DeepEqual(resultArr, expectedArr) {
		t.Errorf("Expected %v but got %v", expectedArr, resultArr)
	}

	expected = "sandwich"
	result = FindDataByAccessKey(mixedData, "favoriteFoods.lunch")
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	// 测试无效的访问键的情况
	expected = []interface{}(nil)
	result = FindDataByAccessKey(mixedData, "pets.*.color")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

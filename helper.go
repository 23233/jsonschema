package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SchemaHelper struct {
	raw        map[string]any
	visited    map[*map[string]any]bool
	accessKeys []string
}

// ResolveRef 解析 JSON schema 中的 $ref 引用，返回引用的 JSON 对象
func (c *SchemaHelper) ResolveRef(ref string) (map[string]any, error) {
	if !strings.HasPrefix(ref, "#") {
		// 不支持非本地引用
		return nil, errors.New("不支持非本地引用")
	}

	parts := strings.Split(strings.TrimPrefix(ref, "#/"), "/")
	target := c.raw
	for _, part := range parts {
		if _, ok := target[part]; !ok {
			return nil, errors.New("未找到对应schema")
		}
		target = target[part].(map[string]any)
	}
	return target, nil
}

func (c *SchemaHelper) SetSchema(input any) error {
	if v, ok := input.(map[string]any); ok {
		c.raw = v
		return nil
	}
	mp, err := StructToMap(input)
	if err != nil {
		return err
	}
	c.raw = mp
	return nil
}

func (c *SchemaHelper) GetRaw() map[string]any {
	return c.raw
}

func (c *SchemaHelper) ToStruct(out any) error {
	return MapToStruct(c.raw, out)
}

func (c *SchemaHelper) GetSchemaMapByPointer(schema map[string]any, pointer string) (map[string]any, error) {
	if len(pointer) < 1 {
		return nil, errors.New("pointer is empty")
	}
	if pointer == "#" || pointer == "/" {
		return schema, nil
	}
	if strings.HasPrefix(pointer, "#") {
		pointer = strings.TrimPrefix(pointer, "#")
	}
	var err error
	parts := strings.Split(strings.TrimPrefix(pointer, "/"), "/")
	for _, part := range parts {
		if part == "" {
			return nil, errors.New("invalid JSON pointer")
		}
		if schema == nil {
			return nil, errors.New("schema is empty")
		}

		if _, ok := schema["type"]; !ok {
			return nil, errors.New("invalid schema type")
		}
		switch schema["type"] {
		case "object":
			properties, ok := schema["properties"].(map[string]any)
			if !ok {
				return nil, errors.New("invalid schema properties")
			}
			if _, ok := properties[part]; !ok {
				return nil, fmt.Errorf("schema properties not has key %s %v", part, schema)
			}
			schema = properties[part].(map[string]any)
			break
		case "array":
			items, ok := schema["items"].(map[string]any)
			if !ok {
				// 那可能items是数组
				itemsArray, ok := schema["items"].([]any)
				if !ok || len(itemsArray) == 0 {
					return nil, errors.New("invalid schema items")
				}

				// 如果items存在是一个数组
				index, err := strconv.Atoi(part)
				if err != nil {
					// 但是传入的序列不是下标数字
					return nil, errors.New("invalid JSON pointer format index")
				}
				// 传入的下标溢出了现有
				if index >= len(itemsArray) {
					return nil, fmt.Errorf("invalid JSON pointer, segment index %d out of range", index)
				}

				// 获取下标的数据不是一个map 那可能是一个空
				arrayItem, ok := itemsArray[index].(map[string]any)
				if !ok {
					return nil, errors.New("invalid JSON pointer, target schema in array but item not a map")
				}

				schema = arrayItem
			} else {
				schema = items
			}
			break

		default:
			return nil, fmt.Errorf("unsupported schema type: %v", schema["type"])
		}
		// 解析ref
		schema, err = c.SchemaRefParse(schema)
		if err != nil {
			return nil, err
		}

	}

	return schema, nil

}

func (c *SchemaHelper) SchemaRefParse(schema map[string]any) (map[string]any, error) {

	// 处理 $ref 引用
	if _, ok := schema["$ref"]; ok {
		ref, ok := schema["$ref"].(string)
		if !ok {
			return nil, errors.New("invalid $ref")
		}

		// 解析引用指向的 schema
		refSchema, err := c.ResolveRef(ref)
		if err != nil {
			return nil, err
		}

		// 如果已经访问过，直接返回
		if c.visited[&refSchema] {
			return nil, fmt.Errorf("circular reference detected in schema: %v", schema)
		}

		// 记录已经访问过的 schema
		c.visited[&refSchema] = true

		// 判断获取出来的ref是否又包含了$ref
		return c.SchemaRefParse(refSchema)
	}
	return schema, nil
}

// 遍历生成accessKey
func (c *SchemaHelper) traverse(currentSchema map[string]any, currentPath string) error {

	schema, err := c.SchemaRefParse(currentSchema)
	if err != nil {
		return err
	}
	typ := schema["type"].(string)
	if typ == "object" {
		if widget, ok := schema["widget"].(string); ok && widget == "RawJsonTree" {
			c.accessKeys = append(c.accessKeys, currentPath)
		} else {

			if properties, ok := schema["properties"].(map[string]any); ok {
				for propertyName, propertySchema := range properties {
					path := propertyName
					if currentPath != "" {
						path = currentPath + "." + propertyName
					}
					_ = c.traverse(propertySchema.(map[string]any), path)
				}
			}
		}

	} else if typ == "array" {
		if items, ok := schema["items"].([]any); ok {
			for index, item := range items {
				path := strconv.Itoa(index)
				if currentPath != "" {
					path = currentPath + "." + path
				}
				_ = c.traverse(item.(map[string]any), path)
			}
		} else if itemsSchema, ok := schema["items"].(map[string]any); ok {

			itemsSchema, err = c.SchemaRefParse(itemsSchema)
			if err != nil {
				return err
			}

			if itemsSchema["type"].(string) == "array" || itemsSchema["type"].(string) == "object" {
				_ = c.traverse(itemsSchema, currentPath+".*")
			} else {
				c.accessKeys = append(c.accessKeys, currentPath)
			}
		}
	} else {
		c.accessKeys = append(c.accessKeys, currentPath)
	}
	return nil
}

// GenAccessKeys 根据json schema生成可访问的accessKey列表
func (c *SchemaHelper) GenAccessKeys() []string {

	if len(c.accessKeys) > 0 {
		return c.accessKeys
	}

	_ = c.traverse(c.raw, "")

	if c.accessKeys[0] == "" {
		c.accessKeys = c.accessKeys[1:]
	}

	return c.accessKeys
}

func NewSchemaHelper(input any) *SchemaHelper {
	var t = new(SchemaHelper)
	_ = t.SetSchema(input)
	t.visited = make(map[*map[string]any]bool)
	t.accessKeys = make([]string, 0)
	return t
}

// GetSchemaMapByPointer 传入一个被序列化之后的 json schema , 和对应需要获取的pointer , 返回 获取到的schema 或者 error
// pointer 格式为 /字段1/字段2 或者 #/字段1/字段2
func GetSchemaMapByPointer(schema map[string]any, pointer string) (map[string]any, error) {
	var t = NewSchemaHelper(schema)
	return t.GetSchemaMapByPointer(t.raw, pointer)
}

func FindDataByAccessKey(data any, accessKey string) any {
	keys := strings.Split(accessKey, ".")
	var currentData = data

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if arrData, ok := currentData.([]any); ok {
			// 处理数组
			if key == "*" {
				// 获取数组的所有元素
				var result []any
				for j := 0; j < len(arrData); j++ {
					elem := FindDataByAccessKey(arrData[j], strings.Join(keys[i+1:], "."))
					if elem != nil {
						switch e := elem.(type) {
						case []any:
							result = append(result, e...)
						default:
							result = append(result, elem)
						}
					}
				}
				return result
			} else if matched, err := regexp.MatchString(`\*\.\d+`, key); err == nil && matched {
				// 获取数组的某个元素
				wildcardIndex := strings.Index(key, "*.")
				index, _ := strconv.Atoi(key[wildcardIndex+2:])
				if index >= len(arrData) {
					return nil
				}
				var result []any
				for j := 0; j < len(arrData); j++ {
					if j == index {
						elem := FindDataByAccessKey(arrData[j], strings.Join(keys[i+2:], "."))
						if elem != nil {
							switch e := elem.(type) {
							case []any:
								result = append(result, e...)
							default:
								result = append(result, elem)
							}
						}
					}
				}
				i += 1
				currentData = result
			} else {
				index, _ := strconv.Atoi(key)

				// 下标超出元素
				if index >= len(arrData) {
					return []any{}
				}

				currentData = arrData[index]
				if currentData == nil {
					return []any{}
				}
			}
		} else if objData, ok := currentData.(map[string]any); ok {
			// 处理对象
			if key == "*" {
				// 获取对象的所有值
				var result []any
				for _, value := range objData {
					elem := FindDataByAccessKey(value, strings.Join(keys[i+1:], "."))
					if elem != nil {
						switch t := elem.(type) {
						case []any:
							result = append(result, t...) // 直接使用类型断言后的变量 t
						default:
							result = append(result, elem)
						}

					}
				}
				return result
			} else if matched, err := regexp.MatchString(`\*\.[^*]+`, key); err == nil && matched {
				// 获取对象的某个值
				wildcardIndex := strings.Index(key, "*.")
				objKey := key[wildcardIndex+2:]
				var result []any
				for _, value := range objData {
					elem := FindDataByAccessKey(value, strings.Join(keys[i+2:], "."))
					if elem != nil {
						switch t := elem.(type) {
						case map[string]any:
							if value, exists := t[objKey]; exists {
								result = append(result, value)
							}
						case []any:
							for _, subElem := range t {
								if subObjElem, ok := subElem.(map[string]any); ok {
									if value, exists := subObjElem[objKey]; exists {
										result = append(result, value)
									}
								}
							}
						}

					}
				}
				i += 1
				currentData = result
			} else {
				currentData = objData[key]
				if currentData == nil {
					return nil
				}
			}
		} else {
			// 不支持其他类型
			return nil
		}
	}
	return currentData
}

// StructToMap 通过json序列化实现struct到map
func StructToMap(in any) (map[string]any, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func MapToStruct[T any](m map[string]any, out T) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}

package jsonschema_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/23233/jsonschema"
)

type SampleUser struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name" jsonschema:"title=the name,description=The name of a friend,example=joe,example=lucy,default=alex"`
	Friends     []int                  `json:"friends,omitempty" jsonschema_description:"The list of IDs, omitted when empty"`
	Tags        map[string]interface{} `json:"tags,omitempty" jsonschema_extras:"a=b,foo=bar,foo=bar1"`
	BirthDate   time.Time              `json:"birth_date,omitempty" jsonschema:"oneof_required=date"`
	YearOfBirth string                 `json:"year_of_birth,omitempty" jsonschema:"oneof_required=year"`
	Metadata    interface{}            `json:"metadata,omitempty" jsonschema:"oneof_type=string;array"`
	FavColor    string                 `json:"fav_color,omitempty" jsonschema:"enum=red,enum=green,enum=blue"`
}

func ExampleReflect() {
	s := jsonschema.Reflect(&SampleUser{})
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(data))
	// Output:
	// {
	//  "$schema": "https://json-schema.org/draft/2020-12/schema",
	//  "$id": "https://github.com/23233/jsonschema_test/sample-user",
	//  "$ref": "#/$defs/SampleUser",
	//  "$defs": {
	//    "SampleUser": {
	//      "oneOf": [
	//        {
	//          "required": [
	//            "birth_date"
	//          ],
	//          "title": "date"
	//        },
	//        {
	//          "required": [
	//            "year_of_birth"
	//          ],
	//          "title": "year"
	//        }
	//      ],
	//      "properties": {
	//        "id": {
	//          "type": "integer",
	//          "attach_data": {
	//            "kind": "int"
	//          }
	//        },
	//        "name": {
	//          "type": "string",
	//          "title": "the name",
	//          "description": "The name of a friend",
	//          "default": "alex",
	//          "examples": [
	//            "joe",
	//            "lucy"
	//          ],
	//          "attach_data": {
	//            "kind": "string"
	//          }
	//        },
	// 		"friends": {
	//          "items": {
	//            "type": "integer",
	//            "attach_data": {
	//              "kind": "int"
	//            }
	//          },
	//          "type": "array",
	//          "description": "The list of IDs, omitted when empty",
	//          "attach_data": {
	//            "kind": "slice"
	//          }
	//        },
	//        "tags": {
	//          "type": "object",
	//          "attach_data": {
	//            "kind": "map"
	//          },
	//          "a": "b",
	//          "foo": [
	//            "bar",
	//            "bar1"
	//          ]
	//        },
	//        "birth_date": {
	//          "type": "string",
	//          "format": "date-time",
	//          "attach_data": {
	//            "kind": "struct"
	//          }
	//        },
	//        "year_of_birth": {
	//          "type": "string",
	//          "attach_data": {
	//            "kind": "string"
	//          }
	//        },
	//        "metadata": {
	//          "oneOf": [
	//            {
	//              "type": "string"
	//            },
	//            {
	//              "type": "array"
	//            }
	//          ],
	//          "attach_data": {
	//            "kind": "interface"
	//          }
	//        },
	//        "fav_color": {
	//          "type": "string",
	//          "enum": [
	//            "red",
	//            "green",
	//            "blue"
	//          ],
	//          "attach_data": {
	//            "kind": "string"
	//          }
	//        }
	//      },
	//      "additionalProperties": false,
	//      "type": "object",
	//      "required": [
	//        "id",
	//        "name"
	//      ],
	//      "attach_data": {
	//        "kind": "struct"
	//      }
	//    }
	//  }
	//}
}

package djson

import (
	"encoding/json"
	"github.com/support-go/session"
	"github.com/support-go/stream/entry"
)

type JsonObjectEntry map[string]interface{}

func (e *JsonObjectEntry) GetValue(key interface{}) (interface{}, session.Status) {
	switch key.(type) {
	case string:
		return (*e)[key.(string)], session.NoError
	}
	return nil, session.JsonArrayKeyTypeErr.ByMessage("读取JsonObject的key类型错误")
}

func (e *JsonObjectEntry) GetIntValue(key interface{}) (int64, session.Status) {
	return getIntValue(e.GetValue(key))
}

func (e *JsonObjectEntry) GetFloatValue(key interface{}) (float64, session.Status) {
	return getFloatValue(e.GetValue(key))
}

type JsonObject struct {
	*entry.Reader
	content JsonObjectEntry
}

func NewJsonObject() *JsonObject {
	return NewJsonObjectFromMap(make(map[string]interface{}))
}

func NewJsonObjectFromMap(input map[string]interface{}) *JsonObject {
	object := &JsonObject{
		content: input,
	}
	object.Reader = entry.NewReader(&object.content)
	return object
}

func JsonObjectFromString(input string) (*JsonObject, error) {
	return JsonObjectFromBuffer([]byte(input))
}

func JsonObjectFromBuffer(input []byte) (*JsonObject, error) {
	object := NewJsonObject()
	err := json.Unmarshal(input, &object.content)
	return object, err
}

func (o *JsonObject) ToJsonString() (string, error) {
	buffer, err := json.Marshal(o.content)
	return string(buffer), err
}

func (o *JsonObject) Foreach(callback func(string, interface{}) bool) {
	for k, v := range o.content {
		if callback(k, v) {
			break
		}
	}
}

func (o *JsonObject) ContainKey(key string) bool {
	_, ok := o.content[key]
	return ok
}

func (o *JsonObject) Put(key string, value interface{}) *JsonObject {
	switch value.(type) {
	case *JsonObject:
		o.content[key] = value.(*JsonObject).content
	case *JsonArray:
		o.content[key] = value.(*JsonArray).content
	default:
		o.content[key] = value
	}
	return o
}

func (o *JsonObject) Length() int {
	return len(o.content)
}

func (o *JsonObject) Unmarshal(ptr interface{}) error {
	buffer, err := json.Marshal(o.content)
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, ptr)
}

func (o *JsonObject) UnmarshalElem(key string, ptr interface{}) error {
	buffer, err := json.Marshal(o.content[key])
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, ptr)
}

func (o *JsonObject) GetJsonObject(key string) *JsonObject {
	value := o.content[key]
	switch value.(type) {
	case map[string]interface{}:
		return NewJsonObjectFromMap(value.(map[string]interface{}))
	}
	return nil
}

func (o *JsonObject) GetJsonArray(key string) *JsonArray {
	value := o.content[key]
	switch value.(type) {
	case []interface{}:
		return NewJsonArrayFromArray(value.([]interface{}))
	}
	return nil
}

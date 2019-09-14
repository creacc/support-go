package djson

import (
	"encoding/json"
	"github.com/support-go/session"
	"github.com/support-go/stream/entry"
)

type JsonArrayEntry []interface{}

func (e *JsonArrayEntry) GetValue(key interface{}) (interface{}, session.Status) {
	var value interface{}
	switch key.(type) {
	case int:
		value = (*e)[key.(int)]
	default:
		return nil, session.JsonArrayKeyTypeErr.ByMessage("读取JsonArray的key类型错误")
	}
	return value, session.NoError
}

func (e *JsonArrayEntry) GetIntValue(key interface{}) (int64, session.Status) {
	return getIntValue(e.GetValue(key))
}

func (e *JsonArrayEntry) GetFloatValue(key interface{}) (float64, session.Status) {
	return getFloatValue(e.GetValue(key))
}

type JsonArray struct {
	*entry.Reader
	content JsonArrayEntry
}

func NewJsonArray() *JsonArray {
	return NewJsonArrayFromArray(make([]interface{}, 0))
}

func NewJsonArrayFromArray(input []interface{}) *JsonArray {
	arr := &JsonArray{
		content: JsonArrayEntry(input),
	}
	arr.Reader = entry.NewReader(&arr.content)
	return arr
}

func JsonArrayFromString(input string) (*JsonArray, error) {
	return JsonArrayFromBuffer([]byte(input))
}

func JsonArrayFromBuffer(input []byte) (*JsonArray, error) {
	array := NewJsonArray()
	err := json.Unmarshal(input, &array.content)
	return array, err
}

func (a *JsonArray) ToJsonString() (string, error) {
	buffer, err := json.Marshal(a.content)
	return string(buffer), err
}

func (a *JsonArray) Put(value interface{}) *JsonArray {
	switch value.(type) {
	case *JsonObject:
		a.content = append(a.content, value.(*JsonObject).content)
	case *JsonArray:
		a.content = append(a.content, value.(*JsonArray).content)
	default:
		a.content = append(a.content, value)
	}
	return a
}

func (a *JsonArray) Foreach(callback func(int, interface{}) bool) {
	for k, v := range a.content {
		if callback(k, v) {
			break
		}
	}
}

func (a *JsonArray) Length() int {
	return len(a.content)
}

func (a *JsonArray) Unmarshal(ptr interface{}) error {
	buffer, err := json.Marshal(a.content)
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, ptr)
}

func (a *JsonArray) UnmarshalElem(index int, ptr interface{}) error {
	buffer, err := json.Marshal(a.content[index])
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, ptr)
}

func (a *JsonArray) GetJsonObject(index int) *JsonObject {
	value := a.content[index]
	switch value.(type) {
	case map[string]interface{}:
		return NewJsonObjectFromMap(value.(map[string]interface{}))
	}
	return nil
}

func (a *JsonArray) GetJsonArray(index int) *JsonArray {
	value := a.content[index]
	switch value.(type) {
	case []interface{}:
		return NewJsonArrayFromArray(value.([]interface{}))
	}
	return nil
}

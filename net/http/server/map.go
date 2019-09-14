package httpsrv

import (
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"mime/multipart"
	"reflect"
)

type MapParams struct {
	m map[string]interface{}
}

func NewMapParams(m map[string]interface{}) *MapParams {
	return &MapParams{
		m: m,
	}
}

func NewStringMapParams(m map[string]string) *MapParams {
	temp := make(map[string]interface{})
	for k, v := range m {
		temp[k] = v
	}
	return NewMapParams(temp)
}

func NewMapParamsFromMap(m map[string]interface{}) *MapParams {
	return &MapParams{
		m: m,
	}
}

func (p *MapParams) getInt(key string) (int64, session.Status) {
	value, stat := p.query(key)
	if stat == session.NoError {
		return formatIntValue(value)
	}
	return 0, stat
}

func (p *MapParams) getFloat(key string) (float64, session.Status) {
	value, stat := p.query(key)
	if stat == session.NoError {
		return formatFloatValue(value)
	}
	return 0, stat
}

func (p *MapParams) getString(key string) (string, session.Status) {
	value, stat := p.query(key)
	if stat == session.NoError {
		return formatStringValue(value)
	}
	return "", stat
}

func (p *MapParams) getFile(key string) (multipart.File, session.Status) {
	return nil, session.FormatErr.ByMessage("转换文件参数格式错误")
}

func (p *MapParams) getParameters(key string) (Parameters, session.Status) {
	value, stat := p.query(key)
	if stat == session.NoError {
		return formatParametersValue(value)
	}
	return nil, stat
}

func (p *MapParams) getArrayParameters(key string) (ArrayParameters, session.Status) {
	value, stat := p.query(key)
	if stat == session.NoError {
		return formatArrayParametersValue(value)
	}
	return nil, stat
}

func (p *MapParams) query(key string) (interface{}, session.Status) {
	value, ok := p.m[key]
	if ok {
		return value, session.NoError
	}
	return nil, session.NotFound.Empty()
}

func formatIntValue(temp interface{}) (int64, session.Status) {
	switch temp.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return temp.(int64), session.NoError
	default:
		log.Logger.Error("Parse query params err : wrong int format %v", temp)
		return 0, session.FormatErr.ByMessage("转换整型参数格式错误")
	}
}

func formatFloatValue(temp interface{}) (float64, session.Status) {
	switch temp.(type) {
	case float32, float64:
		return temp.(float64), session.NoError
	default:
		log.Logger.Error("Parse query params err : wrong float format %v", temp)
		return 0, session.FormatErr.ByMessage("转换浮点型参数格式错误")
	}
}

func formatStringValue(temp interface{}) (string, session.Status) {
	switch temp.(type) {
	case string:
		return temp.(string), session.NoError
	default:
		log.Logger.Error("Parse query params err : wrong string format %v", temp)
		return "", session.FormatErr.ByMessage("转换字符串参数格式错误")
	}
}

func formatParametersValue(temp interface{}) (Parameters, session.Status) {
	vo := reflect.ValueOf(temp)
	if vo.Kind() == reflect.Struct {
		return NewMapParamsFromMap(temp.(map[string]interface{})), session.NoError
	} else {
		log.Logger.Error("Parse query params err : wrong struct format %v", temp)
		return nil, session.FormatErr.ByMessage("转换对象参数格式错误")
	}
}

func formatArrayParametersValue(temp interface{}) (ArrayParameters, session.Status) {
	vo := reflect.ValueOf(temp)
	if vo.Kind() == reflect.Array {
		return NewJsonArrayParams(temp.([]interface{})), session.NoError
	} else {
		log.Logger.Error("Parse query params err : wrong array format %v", temp)
		return nil, session.FormatErr.ByMessage("转换数组参数格式错误")
	}
}

type JsonArrayParams struct {
	array []interface{}
}

func NewJsonArrayParams(array []interface{}) *JsonArrayParams {
	return &JsonArrayParams{
		array: array,
	}
}

func (p *JsonArrayParams) getInt(index int) (int64, session.Status) {
	value, stat := p.query(index)
	if stat.Successful() {
		return formatIntValue(value)
	}
	return 0, stat
}

func (p *JsonArrayParams) getFloat(index int) (float64, session.Status) {
	value, stat := p.query(index)
	if stat.Successful() {
		return formatFloatValue(value)
	}
	return 0, stat
}

func (p *JsonArrayParams) getString(index int) (string, session.Status) {
	value, stat := p.query(index)
	if stat.Successful() {
		return formatStringValue(value)
	}
	return "", stat
}

func (p *JsonArrayParams) getParameters(index int) (Parameters, session.Status) {
	value, stat := p.query(index)
	if stat.Successful() {
		return formatParametersValue(value)
	}
	return nil, stat
}

func (p *JsonArrayParams) getArrayParameters(index int) (ArrayParameters, session.Status) {
	value, stat := p.query(index)
	if stat.Successful() {
		return formatArrayParametersValue(value)
	}
	return nil, stat
}

func (p *JsonArrayParams) query(index int) (interface{}, session.Status) {
	if index >= p.getLength() {
		return nil, session.NotFound.Empty()
	}
	return p.array[index], session.NoError
}

func (p *JsonArrayParams) getLength() int {
	if p.array == nil {
		return 0
	}
	return len(p.array)
}

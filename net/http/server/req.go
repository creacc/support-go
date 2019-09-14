package httpsrv

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"mime/multipart"
	"reflect"
	"strconv"
)

const (
	TagKey      = "key"
	TagRequired = "required"
)

type ArrayParameters interface {
	getInt(index int) (int64, session.Status)
	getFloat(index int) (float64, session.Status)
	getString(index int) (string, session.Status)
	getParameters(index int) (Parameters, session.Status)
	getArrayParameters(index int) (ArrayParameters, session.Status)
	getLength() int
}

type Parameters interface {
	getInt(key string) (int64, session.Status)
	getFloat(key string) (float64, session.Status)
	getString(key string) (string, session.Status)
	getFile(key string) (multipart.File, session.Status)
	getParameters(key string) (Parameters, session.Status)
	getArrayParameters(key string) (ArrayParameters, session.Status)
}

func ParseRequest(params Parameters, callback interface{}) (interface{}, session.Status) {
	to := reflect.TypeOf(callback)
	if to.NumIn() == 1 {
		value, err := unmarshal(params, to.In(0))
		if err.Failed() {
			return nil, err
		}
		values := reflect.ValueOf(callback).Call([]reflect.Value{value.Elem()})
		return values[0].Interface(), session.NoError
	}
	return nil, session.ServerReqErr.ByMessage("服务器错误")
}

func unmarshal(params Parameters, reqType reflect.Type) (reqValue reflect.Value, err session.Status) {
	reqValue = reflect.New(reqType)
	count := reqType.NumField()
	err = session.NoError
	for i := 0; i < count; i++ {
		fieldTo := reqType.Field(i)
		tag := fieldTo.Tag
		if key, ok := tag.Lookup("key"); ok {
			defVal, allowNil := tag.Lookup("default")
			fieldVo := reqValue.Elem().Field(i)
			switch fieldVo.Interface().(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				err = parseInt(fieldVo, key, allowNil, defVal, params)
			case float32, float64:
				err = parseFloat(fieldVo, key, allowNil, defVal, params)
			case string:
				err = parseString(fieldVo, key, allowNil, defVal, params)
			default:
				if fieldVo.Type().String() == "multipart.File" {
					err = parseFile(fieldVo, key, allowNil, params)
				} else {
					switch fieldVo.Kind() {
					case reflect.Struct:
						err = parseStruct(fieldVo, key, allowNil, defVal, params)
					case reflect.Array:
						err = parseArray(fieldVo, key, allowNil, defVal, params)
					default:
						err = session.ServerReqParamTypeUnknownErr.ByMessageFormat("未知参数类型[%s]", key)
					}
				}
			}
			if err.Failed() {
				return
			}
		}
	}
	return
}

func parseInt(fieldVo reflect.Value, key string, allowNil bool, defVal string, params Parameters) (err session.Status) {
	value, status := params.getInt(key)
	switch {
	case status.Successful():
		fieldVo.SetInt(value)
	case status.NotFound():
		if allowNil {
			v, e := strconv.ParseInt(defVal, 10, 64)
			if e != nil {
				err = session.ServerReqLackErr.ByMessageFormat(" 缺省参数[%s]格式错误:%s", key, defVal)
				log.Logger.Error("parse request [%s] err : %v", key, err)
			} else {
				fieldVo.SetInt(v)
			}
		} else {
			err = session.ServerReqLackErr.ByMessageFormat(" 缺少参数[%s]", key)
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	case status.InstanceOf(session.FormatErr):
		err = session.FormatErr.ByMessageFormat(" 参数[%s]格式错误:%s", key, defVal)
		log.Logger.Error("parse request [%s] err : %v", err)
	}
	return
}

func parseFloat(fieldVo reflect.Value, key string, allowNil bool, defVal string, params Parameters) (err session.Status) {
	value, status := params.getFloat(key)
	switch {
	case status.Successful():
		fieldVo.SetFloat(value)
	case status.NotFound():
		if allowNil {
			v, e := strconv.ParseFloat(defVal, 10)
			if e != nil {
				err = session.FormatErr.ByMessageFormat(" 缺省参数[%s]格式错误:%s", key, defVal)
				log.Logger.Error("parse request [%s] err : %v", key, err)
			} else {
				fieldVo.SetFloat(v)
			}
		} else {
			err = session.ServerReqLackErr.ByMessage(fmt.Sprintf(" 缺少参数[%s]", key))
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	case status.InstanceOf(session.FormatErr):
		err = session.FormatErr.ByMessageFormat(" 参数[%s]格式错误:%s", key, defVal)
		log.Logger.Error("parse request [%s] err : %v", err)
	}
	return
}

func parseString(fieldVo reflect.Value, key string, allowNil bool, defVal string, params Parameters) (err session.Status) {
	value, status := params.getString(key)
	switch {
	case status.Successful():
		fieldVo.SetString(value)
	case status.NotFound():
		if allowNil {
			fieldVo.SetString(defVal)
		} else {
			err = session.ServerReqLackErr.ByMessage(fmt.Sprintf(" 缺少参数[%s]", key))
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	}
	return
}

func parseFile(fieldVo reflect.Value, key string, allowNil bool, params Parameters) (err session.Status) {
	value, status := params.getFile(key)
	switch {
	case status.Successful():
		fieldVo.Set(reflect.ValueOf(value))
	case status.NotFound():
		if allowNil == false {
			err = session.ServerReqLackErr.ByMessageFormat(" 缺少参数[%s]", key)
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	}
	return
}

func parseStruct(fieldVo reflect.Value, key string, allowNil bool, defVal string, params Parameters) (err session.Status) {
	params, status := params.getParameters(key)
	switch {
	case status.Successful():
		value, err := unmarshal(params, fieldVo.Type())
		if err.Failed() {
			return err
		}
		fieldVo.Set(value)
	case status.NotFound():
		if allowNil == false {
			err = session.ServerReqLackErr.ByMessage(fmt.Sprintf(" 缺少参数[%s]", key))
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	case status.InstanceOf(session.FormatErr):
		err = session.FormatErr.ByMessageFormat(" 参数[%s]格式错误:%s", key, defVal)
		log.Logger.Error("parse request [%s] err : %v", err)
	}
	return
}

func parseArray(fieldVo reflect.Value, key string, allowNil bool, defVal string, params Parameters) (err session.Status) {
	arrayParams, status := params.getArrayParameters(key)
	switch {
	case status.Successful():
		length := arrayParams.getLength()
		values := make([]reflect.Value, length)
		for i := 0; i < length; i++ {
			params, status = arrayParams.getParameters(i)
			switch {
			case status.Successful():
				value, err := unmarshal(params, fieldVo.Type())
				if err.Failed() {
					return err
				}
				values = append(values, value)
			case status.NotFound():
				if allowNil == false {
					err = session.ServerReqLackErr.ByMessage(fmt.Sprintf(" 缺少参数[%s]", key))
					log.Logger.Error("parse request [%s] err : %v", key, err)
					return
				}
			case status.InstanceOf(session.FormatErr):
				err = session.FormatErr.ByMessageFormat(" 参数[%s]格式错误:%s", key, defVal)
				log.Logger.Error("parse request [%s] err : %v", err)
				return
			}
		}
		fieldVo.Set(reflect.Append(fieldVo, values...))
	case status.NotFound():
		if allowNil == false {
			err = session.ServerReqLackErr.ByMessage(fmt.Sprintf(" 缺少参数[%s]", key))
			log.Logger.Error("parse request [%s] err : %v", key, err)
		}
	case status.InstanceOf(session.FormatErr):
		err = session.FormatErr.ByMessageFormat(" 参数[%s]格式错误:%s", key, defVal)
		log.Logger.Error("parse request [%s] err : %v", err)
	}
	return
}

func QueryIntWithDefault(ctx *gin.Context, key string) int {
	temp := ctx.Query(key)
	if len(temp) == 0 {
		return 0
	}
	value, err := strconv.Atoi(temp)
	if err != nil {
		return 0
	}
	return value
}

func QueryInt(ctx *gin.Context, key string) (int, session.Status) {
	temp := ctx.Query(key)
	if len(temp) == 0 {
		return 0, session.ServerReqLackErr.ByMessageFormat("未找到参数[%s]", key)
	}
	value, err := strconv.Atoi(temp)
	if err != nil {
		return 0, session.FormatErr.NewStatusFormat(err, "参数无法转换为整型[%s]", key)
	}
	return value, session.NoError
}

func QueryInt64WithDefault(ctx *gin.Context, key string) int64 {
	temp := ctx.Query(key)
	if len(temp) == 0 {
		return 0
	}
	value, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func QueryInt64(ctx *gin.Context, key string) (int64, session.Status) {
	temp := ctx.Query(key)
	if len(temp) == 0 {
		return 0, session.ServerReqLackErr.ByMessageFormat("未找到参数[%s]", key)
	}
	value, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		return 0, session.FormatErr.NewStatusFormat(err, "参数无法转换为64位整型[%s]", key)
	}
	return value, session.NoError
}

func QueryFloat64(ctx *gin.Context, key string) (float64, session.Status) {
	temp := ctx.Query(key)
	if len(temp) == 0 {
		return 0, session.ServerReqLackErr.ByMessageFormat("未找到参数[%s]", key)
	}
	value, err := strconv.ParseFloat(temp, 10)
	if err != nil {
		return 0, session.FormatErr.NewStatusFormat(err, "参数无法转换为64位浮点型[%s]", key)
	}
	return value, session.NoError
}

func isRequired(tag reflect.StructTag) bool {
	required, ok := tag.Lookup(TagRequired)
	if ok == false || required == "true" {
		return true
	}
	return false
}

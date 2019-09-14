package djson

import (
	"fmt"
	"github.com/support-go/session"
	"strconv"
)

type JsonElement interface {
	ToJsonString() (string, error)
}

func getIntValue(value interface{}, stat session.Status) (int64, session.Status) {
	if stat.Failed() {
		return 0, stat
	}
	switch value.(type) {
	case float64:
		temp := value.(float64)
		return int64(temp), session.NoError
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32:
		temp := fmt.Sprint(value)
		i, _ := strconv.ParseInt(temp, 10, 64)
		return i, session.NoError
	}
	return 0, session.JsonArrayValueTypeErr.ByMessage("数据类型转换失败")
}

func getFloatValue(value interface{}, stat session.Status) (float64, session.Status) {
	if stat.Failed() {
		return 0, stat
	}
	switch value.(type) {
	case float32:
		return float64(value.(float32)), session.NoError
	case float64:
		return value.(float64), session.NoError
	}
	return 0, session.JsonArrayValueTypeErr.ByMessage("数据类型转换失败")
}

package entry

import (
	"fmt"
	"github.com/support-go/session"
	"strconv"
)

type Entry interface {
	GetValue(key interface{}) (interface{}, session.Status)
	GetIntValue(key interface{}) (int64, session.Status)
	GetFloatValue(key interface{}) (float64, session.Status)
}

type ExceptionHandler func(key, value interface{}, status session.Status) interface{}

type Reader struct {
	entry            Entry
	exceptionHandler ExceptionHandler
}

func NewReader(entry Entry) *Reader {
	return &Reader{
		entry: entry,
	}
}

func (r *Reader) SetExceptionHandler(handler ExceptionHandler) *Reader {
	r.exceptionHandler = handler
	return r
}

func (r *Reader) ThrowException(key, value interface{}, status session.Status) interface{} {
	if r.exceptionHandler != nil {
		return r.exceptionHandler(key, value, status)
	}
	return nil
}

func (r *Reader) GetValue(key interface{}) interface{} {
	value, stat := r.entry.GetValue(key)
	if stat.Failed() {
		return nil
	}
	return value
}

func (r *Reader) GetInt(key interface{}) int {
	return int(r.GetInt64(key))
}

func (r *Reader) GetInt64(key interface{}) int64 {
	value, stat := r.entry.GetIntValue(key)
	if stat.Failed() && r.exceptionHandler != nil {
		if r.exceptionHandler != nil {
			defaultValue := r.exceptionHandler(key, value, stat)
			temp := fmt.Sprintf("%v", defaultValue)
			value, _ = strconv.ParseInt(temp, 10, 64)
		}
	}
	return value
}

func (r *Reader) GetInt32(key interface{}) int32 {
	return int32(r.GetInt64(key))
}

func (r *Reader) GetInt16(key interface{}) int16 {
	return int16(r.GetInt64(key))
}

func (r *Reader) GetInt8(key interface{}) int8 {
	return int8(r.GetInt64(key))
}

func (r *Reader) GetUint(key interface{}) uint {
	return uint(r.GetInt64(key))
}

func (r *Reader) GetUint64(key interface{}) uint64 {
	return uint64(r.GetInt64(key))
}

func (r *Reader) GetUint32(key interface{}) uint32 {
	return uint32(r.GetInt64(key))
}

func (r *Reader) GetUint16(key interface{}) uint16 {
	return uint16(r.GetInt64(key))
}

func (r *Reader) GetUint8(key interface{}) uint8 {
	return uint8(r.GetInt64(key))
}

func (r *Reader) GetByte(key interface{}) byte {
	return byte(r.GetInt64(key))
}

func (r *Reader) GetFloat32(key interface{}) float32 {
	return float32(r.GetFloat64(key))
}

func (r *Reader) GetFloat64(key interface{}) float64 {
	value, stat := r.entry.GetFloatValue(key)
	if stat.Failed() && r.exceptionHandler != nil {
		if r.exceptionHandler != nil {
			defaultValue := r.exceptionHandler(key, value, stat)
			temp := fmt.Sprintf("%v", defaultValue)
			value, _ = strconv.ParseFloat(temp, 64)
		}
	}
	return value
}

func (r *Reader) GetString(key interface{}) string {
	return r.GetValue(key).(string)
}

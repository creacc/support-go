package yml

import (
	"fmt"
	"github.com/support-go/session"
	"github.com/support-go/stream/entry"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

type YAMLEntry map[interface{}]interface{}

func (e *YAMLEntry) GetValue(key interface{}) (interface{}, session.Status) {
	return (*e)[key], session.NoError
}

func (e *YAMLEntry) GetIntValue(key interface{}) (int64, session.Status) {
	value, stat := e.GetValue(key)
	if stat.Failed() {
		return 0, stat
	}
	temp := fmt.Sprint(value)
	i, _ := strconv.ParseInt(temp, 10, 64)
	return i, session.NoError
}

func (e *YAMLEntry) GetFloatValue(key interface{}) (float64, session.Status) {
	value, stat := e.GetValue(key)
	if stat.Failed() {
		return 0, stat
	}
	temp := fmt.Sprint(value)
	i, _ := strconv.ParseFloat(temp, 64)
	return i, session.NoError
}

type Reader struct {
	*entry.Reader
	entry YAMLEntry
}

func NewReader(path string) (*Reader, session.Status) {
	reader := &Reader{
		entry: make(map[interface{}]interface{}),
	}
	reader.Reader = entry.NewReader(&reader.entry)
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, session.YMLReadFileErr.NewStatus(err, "读取yml文件失败")
	}
	err = yaml.Unmarshal(buf, &reader.entry)
	if err != nil {
		return nil, session.YMLReadFileErr.NewStatus(err, "解析yml文件失败")
	}
	return reader, session.NoError
}

func (r *Reader) GetChild(key interface{}) (*Reader, session.Status) {
	value, stat := r.entry.GetValue(key)
	if stat.Failed() {
		return nil, stat
	}
	m := value.(YAMLEntry)
	reader := &Reader{
		entry: m,
	}
	reader.Reader = entry.NewReader(&reader.entry)
	return reader, session.NoError
}

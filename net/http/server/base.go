package httpsrv

import (
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"mime/multipart"
	"strconv"
)

type Adapter interface {
	query(key string) (string, session.Status)
}

type BaseParams struct {
	adapter Adapter
}

func (p *BaseParams) getInt(key string) (int64, session.Status) {
	temp, stat := p.adapter.query(key)
	if stat.Successful() {
		return formatInt(temp)
	}
	return 0, stat
}

func (p *BaseParams) getFloat(key string) (float64, session.Status) {
	temp, stat := p.adapter.query(key)
	if stat.Successful() {
		return formatFloat(temp)
	}
	return 0, stat
}

func (p *BaseParams) getString(key string) (string, session.Status) {
	temp, stat := p.adapter.query(key)
	if stat.Successful() {
		return temp, session.NoError
	}
	return "", stat
}

func (p *BaseParams) getFile(key string) (multipart.File, session.Status) {
	return nil, session.FormatErr.ByMessage("获取文件参数格式错误")
}

func (p *BaseParams) getParameters(key string) (Parameters, session.Status) {
	return nil, session.FormatErr.ByMessage("获取对象参数格式错误")
}

func (p *BaseParams) getArrayParameters(key string) (ArrayParameters, session.Status) {
	return nil, session.FormatErr.ByMessage("获取数组参数格式错误")
}

func formatInt(temp string) (int64, session.Status) {
	value, err := strconv.ParseInt(temp, 10, 64)
	if err != nil {
		log.Logger.Error("Parse query params err : %v", err)
		return 0, session.FormatErr.ByMessage("转换整型参数格式错误")
	} else {
		return value, session.NoError
	}
}

func formatFloat(temp string) (float64, session.Status) {
	value, err := strconv.ParseFloat(temp, 10)
	if err != nil {
		log.Logger.Error("Parse query params err : %v", err)
		return 0, session.FormatErr.ByMessage("转换浮点型参数格式错误")
	} else {
		return value, session.NoError
	}
}

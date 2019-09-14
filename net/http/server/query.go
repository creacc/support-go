package httpsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"strings"
)

type QueryParams struct {
	BaseParams
	ctx *gin.Context
}

func (p *QueryParams) query(key string) (string, session.Status) {
	value := p.ctx.Query(key)
	if len(value) == 0 {
		return "", session.NotFound.Empty()
	}
	return value, session.NoError
}

func NewQueryParams(ctx *gin.Context) *QueryParams {
	qp := &QueryParams{
		ctx: ctx,
	}
	qp.adapter = qp
	return qp
}

func (p *QueryParams) getArrayParameters(key string) (ArrayParameters, session.Status) {
	temp, stat := p.query(key)
	if stat.Successful() {
		return formatArrayParameters(temp)
	}
	return nil, stat
}

func formatArrayParameters(temp string) (ArrayParameters, session.Status) {
	length := len(temp)
	if length == 0 || temp[0] != '[' || temp[length-1] != ']' {
		log.Logger.Error("Parse query params err : wrong array format")
		return nil, session.FormatErr.ByMessage("获取数组参数格式错误")
	}
	return NewQueryArrayParams(strings.Split(temp[1:length-1], ",")), session.NoError
}

type QueryArrayParams struct {
	values []string
}

func NewQueryArrayParams(values []string) *QueryArrayParams {
	return &QueryArrayParams{
		values: values,
	}
}
func (p *QueryArrayParams) getInt(index int) (int64, session.Status) {
	temp, stat := p.query(index)
	if stat.Successful() {
		return formatInt(temp)
	}
	return 0, stat
}

func (p *QueryArrayParams) getFloat(index int) (float64, session.Status) {
	temp, stat := p.query(index)
	if stat.Successful() {
		return formatFloat(temp)
	}
	return 0, stat
}

func (p *QueryArrayParams) getString(index int) (string, session.Status) {
	temp, stat := p.query(index)
	if stat.Successful() {
		return temp, session.NoError
	}
	return "", stat
}

func (p *QueryArrayParams) getParameters(index int) (Parameters, session.Status) {
	return nil, session.FormatErr.ByMessage("获取对象参数格式错误")
}

func (p *QueryArrayParams) getArrayParameters(index int) (ArrayParameters, session.Status) {
	temp, stat := p.query(index)
	if stat.Successful() {
		return formatArrayParameters(temp)
	}
	return nil, stat
}

func (p *QueryArrayParams) getLength() int {
	if p.values == nil {
		return 0
	}
	return len(p.values)
}

func (p *QueryArrayParams) query(index int) (string, session.Status) {
	if index >= p.getLength() {
		return "", session.NotFound.Empty()
	}
	return p.values[index], session.NoError
}

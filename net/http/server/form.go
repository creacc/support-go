package httpsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/support-go/session"
)

type FormParams struct {
	BaseParams
	ctx *gin.Context
}

func (p *FormParams) query(key string) (string, session.Status) {
	value := p.ctx.PostForm(key)
	if len(value) == 0 {
		return "", session.NotFound.Empty()
	}
	return value, session.NoError
}

func NewFormParams(ctx *gin.Context) *FormParams {
	fp := &FormParams{
		ctx: ctx,
	}
	fp.adapter = fp
	return fp
}

func (p *FormParams) getParameters(key string) (Parameters, session.Status) {
	m := p.ctx.PostFormMap(key)
	if len(m) == 0 {
		return nil, session.NotFound.Empty()
	}
	return NewStringMapParams(m), session.NoError
}

func (p *FormParams) getArrayParameters(key string) (ArrayParameters, session.Status) {
	a := p.ctx.PostFormArray(key)
	if len(a) == 0 {
		return nil, session.NotFound.Empty()
	}
	return NewQueryArrayParams(a), session.NoError
}

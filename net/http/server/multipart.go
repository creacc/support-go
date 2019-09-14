package httpsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"mime/multipart"
)

type MultiFormParams struct {
	BaseParams
	values map[string][]string
	files  map[string][]*multipart.FileHeader
}

func NewMultiFormParams(ctx *gin.Context) (*MultiFormParams, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Logger.Error("NewMultiFormParams error %v", err)
		return nil, err
	}
	mp := &MultiFormParams{
		values: form.Value,
		files:  form.File,
	}
	mp.adapter = mp
	return mp, nil
}

func (p *MultiFormParams) query(key string) (string, session.Status) {
	value, ok := p.values[key]
	if ok == false {
		return "", session.NotFound.Empty()
	}
	if len(value) == 0 {
		return "", session.NotFound.Empty()
	}
	return value[0], session.NoError
}

func (p *MultiFormParams) getFile(key string) (multipart.File, session.Status) {
	value, ok := p.files[key]
	if ok == false {
		return nil, session.NotFound.Empty()
	}
	if len(value) == 0 {
		return nil, session.NotFound.Empty()
	}
	file, err := value[0].Open()
	if err != nil {
		log.Logger.Error("getFile error %v", err)
		return nil, session.ServerReqFileOpenErr.ByMessage("转换数组参数格式错误")
	}
	return file, session.NoError
}

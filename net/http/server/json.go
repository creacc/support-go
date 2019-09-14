package httpsrv

import "github.com/gin-gonic/gin"

func NewJsonParams(ctx *gin.Context) (params *MapParams, err error) {
	m := make(map[string]interface{})
	err = ctx.BindJSON(m)
	if err == nil {
		params = NewMapParams(m)
	}
	return
}

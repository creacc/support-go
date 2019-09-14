package httpcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"io"
	"mime/multipart"
	"net/url"
	"os"
)

func getQueryParams(req *Request) string {
	params := ""
	for k, v := range req.params {
		params += fmt.Sprintf("&%s=%v", k, v)
	}
	if len(params) > 0 {
		r := []rune(params)
		r[0] = '?'
		params = string(r)
	}
	log.Logger.Debug("[HTTP_C] params : %s", params)
	return params
}

func getJsonParams(req *Request) ([]byte, session.Status) {
	j, err := json.Marshal(req.params)
	if err != nil {
		if err != nil {
			log.Logger.Error("getJsonParams err %v", err)
			return nil, session.ClientReqJsonParamErr.NewStatus(err, "请求失败")
		}
	}
	log.Logger.Debug("[HTTP_C] params : %s", string(j))
	return j, session.NoError
}

func getFormParams(req *Request) url.Values {
	values := url.Values{}
	for k, v := range req.params {
		values.Set(k, fmt.Sprint(v))
	}
	log.Logger.Debug("[HTTP_C] params : %v", values)
	return values
}

func getMultipartFormParams(req *Request) (io.Reader, session.Status) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	defer writer.Close()
	req.SetHeader("Content-Type", writer.FormDataContentType())
	for k, v := range req.params {
		err := writer.WriteField(k, fmt.Sprint(v))
		if err != nil {
			log.Logger.Error("http post param err %v", err)
			return nil, session.ClientReqMultipartParamFieldErr.NewStatus(err, "请求失败")
		}
	}
	for k, v := range req.files {
		name := v.name
		reader := v.reader
		if reader == nil {
			f, err := os.Open(v.name)
			if err != nil {
				log.Logger.Error("http post file open err %v", err)
				return nil, session.ClientReqFileParamFileOpenErr.NewStatus(err, "请求失败")
			}
			reader = f
			name = f.Name()
		}
		fw, err := writer.CreateFormFile(k, name)
		if err != nil {
			log.Logger.Error("http post file err %v", err)
			return nil, session.ClientReqFileParamFileCreateErr.NewStatus(err, "请求失败")
		}
		l, err := io.Copy(fw, reader)
		if err != nil {
			log.Logger.Error("http post file copy err %v", err)
			return nil, session.ClientReqFileParamFileCopyErr.NewStatus(err, "请求失败")
		}
		if l == 0 {
			log.Logger.Error("http post file empty err %v", err)
			return nil, session.ClientReqFileParamFileEmptyErr.NewStatus(err, "请求失败")
		}
	}
	log.Logger.Debug("[HTTP_C] params : %v", req.params)
	log.Logger.Debug("[HTTP_C] files : %v", req.files)
	return &buf, session.NoError
}

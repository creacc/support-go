package httpcli

import (
	"bytes"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"io"
	"net/http"
	"strings"
)

type executor func(req *Request) *Response

var executorMap = map[string]executor{
	http.MethodDelete: func(req *Request) *Response {
		return delete(req)
	},
	http.MethodGet: func(req *Request) *Response {
		return get(req)
	},
	http.MethodPost: func(req *Request) *Response {
		return post(req)
	},
}

func get(req *Request) *Response {
	req.url += getQueryParams(req)
	return doRequest(req, nil)
}

func delete(req *Request) *Response {
	req.url += getQueryParams(req)
	return doRequest(req, nil)
}

func post(req *Request) *Response {
	if req.isJson {
		j, stat := getJsonParams(req)
		if stat.Failed() {
			return newResponse(nil, stat)
		}
		req.headers["Content-Type"] = "application/json"
		return doRequest(req, bytes.NewBuffer(j))
	} else {
		if len(req.files) == 0 {
			values := getFormParams(req)
			req.headers["Content-Type"] = "application/x-www-form-urlencoded"
			return doRequest(req, strings.NewReader(values.Encode()))
		} else {
			buf, stat := getMultipartFormParams(req)
			if stat.Failed() {
				return newResponse(nil, stat)
			}
			return doRequest(req, buf)
		}
	}
}

func doRequest(req *Request, body io.Reader) *Response {
	preprocess(req.preprocessors, req)
	log.Logger.Debug("[HTTP_C] %s", req.url)
	log.Logger.Debug("[HTTP_C] %s", req.method)
	log.Logger.Debug("[HTTP_C] %v", req.headers)
	var stat session.Status
	var httpResp *http.Response
	httpReq, httpErr := http.NewRequest(req.method, req.url, body)
	if httpErr != nil {
		log.Logger.Debug("[HTTP_C] New request failed for %v", httpErr)
		stat = session.ClientReqDoErr.NewStatus(httpErr, "请求失败")
	} else {
		for k, v := range req.headers {
			httpReq.Header.Add(k, v)
		}
		httpResp, httpErr = http.DefaultClient.Do(httpReq)
		if httpErr != nil {
			log.Logger.Debug("[HTTP_C] New request failed for %v", httpErr)
			stat = session.ClientReqDoErr.NewStatus(httpErr, "请求失败")
		} else {
			stat = session.NoError
		}
	}
	resp := newResponse(httpResp, stat)
	process(req.processors, resp)
	return resp
}

func preprocess(preprocessors []Preprocessor, req *Request) {
	if preprocessors != nil {
		for _, v := range preprocessors {
			v(req)
		}
	}
}

func process(processors []Processor, resp *Response) {
	if processors != nil {
		for _, v := range processors {
			v(resp)
		}
	}
}

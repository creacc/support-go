package httpcli

import (
	"github.com/support-go/utils/log"
	"io"
	"net/http"
)

type Preprocessor func(req *Request)

type Processor func(resp *Response)

type Request struct {
	url           string
	method        string
	headers       map[string]string
	params        map[string]interface{}
	files         map[string]FilePart
	preprocessors []Preprocessor
	processors    []Processor
	isJson        bool
}

type FilePart struct {
	name   string
	reader io.Reader
}

func newDefaultRequest(url, method string) *Request {
	return &Request{
		url:     url,
		method:  method,
		headers: make(map[string]string),
		params:  make(map[string]interface{}),
	}
}

func Get(url string) *Request {
	return newDefaultRequest(url, http.MethodGet)
}

func Delete(url string) *Request {
	return newDefaultRequest(url, http.MethodDelete)
}

func Post(url string) *Request {
	r := newDefaultRequest(url, http.MethodPost)
	r.files = make(map[string]FilePart)
	return r
}

func PostJson(url string) *Request {
	r := newDefaultRequest(url, http.MethodPost)
	r.isJson = true
	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

func (r *Request) SetParam(key string, value interface{}) *Request {
	r.params[key] = value
	return r
}

func (r *Request) SetParams(m map[string]interface{}) *Request {
	for k, v := range m {
		r.SetParam(k, v)
	}
	return r
}

func (r *Request) AddPreprocessor(p Preprocessor) *Request {
	r.preprocessors = append(r.preprocessors, p)
	return r
}

func (r *Request) AddProcessor(p Processor) *Request {
	r.processors = append(r.processors, p)
	return r
}

func (r *Request) Params() map[string]interface{} {
	return r.params
}

func (r *Request) IsJsonFormat() bool {
	return r.isJson
}

func (r *Request) SetLocalFile(key, path string) *Request {
	if r.files == nil {
		panic("file param is not support")
	}
	r.files[key] = FilePart{name: path}
	return r
}

func (r *Request) SetFile(key, name string, reader io.Reader) *Request {
	if r.files == nil {
		panic("file param is not support")
	}
	r.files[key] = FilePart{name: name, reader: reader}
	return r
}

func (r *Request) Do() *Response {
	resp := executorMap[r.method](r)
	if resp.Status().Failed() {
		log.Logger.Error("[HTTP_C] Response err : %v", resp.Status())
	}
	return resp
}

func (r *Request) DoAsync(handler func(resp *Response)) {
	go func() {
		resp := executorMap[r.method](r)
		if resp.Status().Failed() {
			log.Logger.Error("[HTTP_C] Response err : %v", resp.Status())
		}
		handler(resp)
	}()
}

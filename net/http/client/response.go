package httpcli

import (
	"encoding/json"
	"github.com/support-go/session"
	"github.com/support-go/utils/log"
	"io/ioutil"
	"net/http"
)

type Response struct {
	resp *http.Response
	stat session.Status
}

func newResponse(resp *http.Response, stat session.Status) *Response {
	return &Response{
		resp: resp,
		stat: stat,
	}
}

func (r *Response) Successful() bool {
	return r.stat.Successful()
}

func (r *Response) Failed() bool {
	return r.stat.Failed()
}

func (r *Response) Status() session.Status {
	return r.stat
}

func (r *Response) StatusCode() int {
	return r.response().StatusCode
}

func (r *Response) String() (string, error) {
	buffer, err := r.ByteArray()
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func (r *Response) ByteArray() ([]byte, error) {
	buffer, err := ioutil.ReadAll(r.response().Body)
	if err != nil {
		log.Logger.Error("Response body err %v", err)
		return nil, err
	}
	return buffer, nil
}

func (r *Response) Json(v interface{}) error {
	buffer, err := r.ByteArray()
	if err != nil {
		log.Logger.Error("Response body json %s err %v", string(buffer), err)
		return err
	}
	return json.Unmarshal(buffer, v)
}

func (r *Response) response() *http.Response {
	if r.resp == nil {
		panic("http resp is nil")
	}
	return r.resp
}

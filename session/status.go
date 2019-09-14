package session

import (
	"fmt"
	"sync"
)

type StatusModule struct {
	sync.Mutex
	current int
	code    int
	desc    string
	offset  int
}

type StatusFactory struct {
	code  int
	desc  string
	empty Status
}

type Status struct {
	code    int
	message string
	err     error
}

var root = &StatusModule{}
var moduleMap = make(map[int]string)
var factoryMap = make(map[int]string)

func (m *StatusModule) shift(code int, desc string, offset int) *StatusModule {
	m.Lock()
	defer m.Unlock()
	if code > m.offset {
		panic("code must less than module offset")
	}
	moduleMap[code] = desc
	return &StatusModule{
		code:   code + m.code*m.offset,
		desc:   desc,
		offset: offset,
	}
}

func (m StatusModule) factory(desc string) *StatusFactory {
	m.Lock()
	defer m.Unlock()
	m.current++
	factoryMap[m.current] = desc
	code := m.current + m.code*m.offset
	return &StatusFactory{
		code:  code,
		desc:  desc,
		empty: Status{code: code},
	}
}

func (f *StatusFactory) ByMessage(msg string) Status {
	return f.NewStatus(nil, msg)
}

func (f *StatusFactory) ByMessageFormat(format string, args ...interface{}) Status {
	return f.ByMessage(fmt.Sprintf(format, args...))
}

func (f *StatusFactory) ByMessageArgs(args ...interface{}) Status {
	return f.ByMessage(fmt.Sprint(args...))
}

func (f *StatusFactory) ByError(err error) Status {
	return f.NewStatus(err, "")
}

func (f *StatusFactory) NewStatus(err error, msg string) Status {
	return Status{
		code:    f.code,
		message: msg,
		err:     err,
	}
}

func (f *StatusFactory) Empty() Status {
	return f.empty
}

func (f *StatusFactory) NewStatusFormat(err error, format string, args ...interface{}) Status {
	return f.NewStatus(err, fmt.Sprintf(format, args...))
}

func (s Status) InstanceOf(factory *StatusFactory) bool {
	return s.code == factory.code
}

func (s Status) Successful() bool {
	return s.code == NoError.code
}

func (s Status) Failed() bool {
	return s.code != NoError.code
}

func (s Status) NotFound() bool {
	return s.code == NotFound.empty.code
}

func (s Status) Error() string {
	return fmt.Sprintf("[%d]%s(%v)", s.code, s.message, s.err)
}

func (s Status) Message() string {
	return s.message
}

func findModuleDesc(code int) string {
	if desc, ok := moduleMap[code]; ok {
		return desc
	}
	return ""
}

func findStatusDesc(code int) string {
	if desc, ok := factoryMap[code]; ok {
		return desc
	}
	return ""
}

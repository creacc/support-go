package session

import uuid "github.com/satori/go.uuid"

type Context struct {
	id     string
	status Status
}

func NewContext() *Context {
	id, _ := uuid.NewV4()
	return &Context{
		id:     id.String(),
		status: Status{},
	}
}

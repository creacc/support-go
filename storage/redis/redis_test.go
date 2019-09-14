package codis

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("127.0.0.1:6379")
	r := c.KeyReader()

	reply, err := r.GetString("aaa")
	fmt.Println(reply)
	fmt.Println(err)
}

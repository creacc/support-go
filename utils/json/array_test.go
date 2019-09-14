package djson

import (
	"fmt"
	"testing"
)

func TestNewJsonArray(t *testing.T) {
	arr := NewJsonArray()
	arr.Put(111)
	arr.Put(222)
	arr.Put(333)

	fmt.Println(arr.GetInt64(1))
}

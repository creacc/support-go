package djson

import (
	"fmt"
	"testing"
)

func TestNewJsonObject(t *testing.T) {
	object, _ := JsonObjectFromString("{\"test\":2.222}")
	//object.Put("aaa", 111)
	//object.Put("bbb", 222)
	//object.Put("ccc", 333)

	fmt.Println(object.GetFloat64("test"))
}

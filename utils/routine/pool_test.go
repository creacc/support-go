package droutine

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	for i := 0; i < 10; i++ {
		DefaultPool.Execute(func() {
			time.Sleep(time.Second)
		})
		fmt.Println("run ", i)
	}
}

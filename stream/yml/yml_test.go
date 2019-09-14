package yml

import (
	"fmt"
	"testing"
)

func TestNewReader(t *testing.T) {
	reader, stat := NewReader("./dev.yml")
	if stat.Failed() {
		fmt.Println(stat.Error())
	}
	mysql, stat := reader.GetChild("mysql")
	if stat.Failed() {
		fmt.Println(stat.Error())
	}
	arthur, stat := mysql.GetChild("gaia")
	if stat.Failed() {
		fmt.Println(stat.Error())
	}
	fmt.Println(arthur.GetString("port"))
}

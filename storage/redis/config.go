package codis

import "fmt"

type Config struct {
	IP   string
	Port string
}

func (c Config) getAddress() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}

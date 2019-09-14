package mysql

import (
	"fmt"
)

type Config struct {
	Username string
	Password string
	IP       string
	Port     string
}

func (c Config) getDataSource(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.IP, c.Port, database)
}

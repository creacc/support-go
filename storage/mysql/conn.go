package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/support-go/utils/log"
)

type Conn struct {
	conf Config
}

func NewConn(conf Config) *Conn {
	return &Conn{
		conf: conf,
	}
}

func (c *Conn) connectToDatabase(database string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.conf.getDataSource(database))
	if err != nil {
		log.Logger.Debug("connect failed %v", err)
		return nil, err
	}
	//db.SetMaxIdleConns(500)  //SetMaxIdleConns用于设置闲置的连接数。
	//db.SetMaxOpenConns(1000) //SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制。
	err = db.Ping()
	if err != nil {
		log.Logger.Debug("connect failed %v", err)
		return nil, err
	}
	log.Logger.Debug("connect success")
	return db, err
}

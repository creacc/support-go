package mysql

import (
	"database/sql"
	"github.com/support-go/utils/log"
)

type Driver interface {
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Context interface {
	Query(query string, args ...interface{}) *QueryResult
	Exec(query string, args ...interface{}) *ExecResult
	close() error
}

type SqlContext struct {
	d      Driver
	failed bool
}

type DBContext struct {
	*SqlContext
	db *sql.DB
}

type TXContext struct {
	*SqlContext
	tx *sql.Tx
}

func wrap(driver Driver) *SqlContext {
	return &SqlContext{
		d: driver,
	}
}

func WrapDB(db *sql.DB) *DBContext {
	return &DBContext{
		SqlContext: wrap(db),
		db:         db,
	}
}

func WrapTX(tx *sql.Tx) *TXContext {
	return &TXContext{
		SqlContext: wrap(tx),
		tx:         tx,
	}
}

func (c *SqlContext) Query(query string, args ...interface{}) *QueryResult {
	log.Logger.Debug("Query SQL: %s", query)
	log.Logger.Debug("Query Args: %v", args)
	stmt, err := c.d.Prepare(query)
	if err != nil {
		c.failed = true
		log.Logger.Error("Query Err: %s in %s", err, query)
		return queryResult(nil, err)
	}
	return queryResult(stmt.Query(args...))
}

func (c *SqlContext) Exec(query string, args ...interface{}) *ExecResult {
	log.Logger.Debug("Exec SQL: %s", query)
	log.Logger.Debug("Exec Args: %v", args)
	stmt, err := c.d.Prepare(query)
	if err != nil {
		c.failed = true
		log.Logger.Error("Prepare Err: %s in %s", err, query)
		return execResult(nil, err)
	}
	return execResult(stmt.Exec(args...))
}

func (c *DBContext) close() error {
	return c.db.Close()
}

func (c *TXContext) close() error {
	if c.failed {
		return c.tx.Rollback()
	} else {
		return c.tx.Commit()
	}
}

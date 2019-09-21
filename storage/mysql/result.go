package mysql

import "database/sql"

type RowHandler func(r *Row)

type QueryResult struct {
	rows    *sql.Rows
	scanned int
	err     error
}

func queryResult(rows *sql.Rows, err error) *QueryResult {
	return &QueryResult{
		rows:    rows,
		scanned: 0,
		err:     err,
	}
}

func (r *QueryResult) Read(handler RowHandler) error {
	if r.err != nil {
		return r.err
	}
	row := &Row{result: r}
	for r.rows.Next() {
		handler(row)
		if r.err != nil {
			return r.err
		}
		r.scanned++
	}
	return nil
}

func (r *QueryResult) Scanned() int {
	return r.scanned
}

func (r *QueryResult) scan(args ...interface{}) {
	r.err = r.rows.Scan(args...)
}

type Row struct {
	result *QueryResult
}

func (r *Row) Scan(args ...interface{}) {
	r.result.scan(args...)
}

type ExecResult struct {
	result sql.Result
	err    error
}

func execResult(result sql.Result, err error) *ExecResult {
	return &ExecResult{
		result: result,
		err:    err,
	}
}

func (r *ExecResult) Id() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.result.LastInsertId()
}

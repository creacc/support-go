package mysql

import "database/sql"

type RowHandler func(r *Row) error

type Rows struct {
	rows *sql.Rows
}

func (r *Rows) Read(handler RowHandler) error {
	row := &Row{rows: r.rows}
	for r.rows.Next() {
		if err := handler(row); err != nil {
			return err
		}
	}
	return nil
}

type Row struct {
	rows *sql.Rows
}

func (r *Row) Scan(args ...interface{}) error {
	return r.rows.Scan(args...)
}

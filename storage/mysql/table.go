package mysql

import "fmt"

const (
	leftJoin  = "left join"
	RightJoin = "right join"
	InnerJoin = "inner join"
)

type Table struct {
	name  string
	alias string
}

type Statement struct {
	db        *Database
	table     *Table
	columns   []string
	condition string
	args      []interface{}
	joiners   []JoinStatement
}

type JoinStatement struct {
	*Statement
	joinType      string
	joinCondition string
}

func (t *Table) Prepare(db *Database, columns ...string) *Statement {
	if columns == nil {
		panic("columns cannot be nil")
	}
	return &Statement{
		db:      db,
		table:   t,
		columns: columns,
	}
}

func (s *Statement) LeftJoin(joiner *Statement, joinColumn, column string) *Statement {
	return s.join(leftJoin, joiner, joinColumn, column)
}

func (s *Statement) RightJoin(joiner *Statement, joinColumn, column string) *Statement {
	return s.join(RightJoin, joiner, joinColumn, column)
}

func (s *Statement) InnerJoin(joiner *Statement, joinColumn, column string) *Statement {
	return s.join(InnerJoin, joiner, joinColumn, column)
}

func (s *Statement) join(joinType string, joiner *Statement, joinColumn, column string) *Statement {
	s.joiners = append(s.joiners, JoinStatement{
		Statement:     joiner,
		joinType:      joinType,
		joinCondition: fmt.Sprintf("%s.%s = %s.%s", joiner.table.alias, joinColumn, s.table.alias, column),
	})
	return s
}

func (s *Statement) Where(condition string, args ...interface{}) *Statement {
	s.condition = condition
	s.args = args
	return s
}

func (s *Statement) Query(handler RowHandler) error {
	return s.db.Run(func(ctx Context) error {
		res := ctx.Query(combineQuery(s), s.args...)
		return res.Read(func(r *Row) {
			handler(r)
		})
	})
}

func combineQuery(s *Statement) string {
	columns := ""
	for i, c := range s.columns {
		if i > 0 {
			columns += ","
		}
		columns += "`" + c + "`"
	}

	args := make([]interface{}, 0)
	join := ""
	if s.joiners != nil {
		for _, j := range s.joiners {
			join += fmt.Sprintf(" %s %s %s on %s", j.joinType, j.table.name, j.table.alias, j.joinCondition)
			if len(j.condition) > 0 {
				join += " and " + j.condition
				args = append(args, j.args)
			}
		}
	}
	return fmt.Sprintf("select %s from %s%s%s", columns, s.table.name, join, s.condition)
}

func combineInsert(s *Statement) string {
	columns := ""
	holders := ""
	for i, c := range s.columns {
		if i > 0 {
			columns += ","
			holders += ","
		}
		columns += "`" + c + "`"
		holders += "?"
	}
	return fmt.Sprintf("insert into %s (%s) values (%s)", s.table.name, columns, holders)
}

func combineUpdate(s *Statement) string {
	columns := ""
	for i, c := range s.columns {
		if i > 0 {
			columns += ","
		}
		columns += "`" + c + "` = ?"
	}
	return fmt.Sprintf("update %s set %s%s", s.table.name, columns, s.condition)
}

func combineDelete(s *Statement) string {
	return fmt.Sprintf("delete from %s%s", s.table.name, s.condition)
}

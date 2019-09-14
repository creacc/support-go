package mysql

type Database struct {
	name string
	conn *Conn
}

type DBTask func(ctx Context) error

func NewDatabase(name string, conn *Conn) *Database {
	return &Database{
		name: name,
		conn: conn,
	}
}

func (d *Database) Run(task DBTask) error {
	db, err := d.conn.connectToDatabase(d.name)
	if err != nil {
		return err
	}
	ctx := WrapDB(db)
	defer ctx.close()
	return task(ctx)
}

func (d *Database) RunWithTransaction(tasks ...DBTask) error {
	db, err := d.conn.connectToDatabase(d.name)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	ctx := WrapTX(tx)
	defer ctx.close()
	for _, v := range tasks {
		err = v(ctx)
		if err != nil {
			break
		}
	}
	return err
}

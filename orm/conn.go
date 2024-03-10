package orm

import "database/sql"

type Connection interface {
	Begin(opts ...*sql.TxOptions) Connection
	Rollback()
	Commit()
	SavePoint(name string)
	RollbackTo(name string)
}

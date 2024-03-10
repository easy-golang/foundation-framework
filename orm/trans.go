package orm

type CommitListener = func()
type RollbackListener = func(err error)
type Transaction interface {
	Do(fun func(conn Connection) error) error

	DoWithListener(fun func(conn Connection) error, commit CommitListener, rollback RollbackListener) error
}

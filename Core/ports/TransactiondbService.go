package ports

type Transaction interface{
	Commit()error
	Rollback()error
}

type TransactionManager interface{
	Begin()(Transaction,error)
}
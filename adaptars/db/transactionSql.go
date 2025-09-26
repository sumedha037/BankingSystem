package db 

import "BankingSystem/Core/ports"
import "database/sql"

func (a *AccountSqlDB) Begin()(ports.Transaction,error){
	tx,err:=a.db.Begin();
	if err!=nil{
		 return nil,err
	}
	return &SqlDBTransaction{Tx:tx},nil
}


type SqlDBTransaction struct{
	Tx  *sql.Tx
}

func (s *SqlDBTransaction)Commit()error{
    return s.Tx.Commit()
}

func (s *SqlDBTransaction)Rollback()error{
	return s.Tx.Rollback()
}


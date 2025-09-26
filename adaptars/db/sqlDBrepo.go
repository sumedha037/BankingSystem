package db

import (
	"BankingSystem/Core/domain"
	 "BankingSystem/customErrors"
	"database/sql"
	"fmt"
	"log"
	"strings"
)


type CustomerSqlDB struct{
	db *sql.DB
}

func NewCustomer(db *sql.DB)*CustomerSqlDB{
   return &CustomerSqlDB{db : db}
}

func (c *CustomerSqlDB)SaveCustomer(customer domain.Customer)error{
	_,err:=c.db.Exec("INSERT INTO Customer (CustomerId , Name , Email , Phone , AccountType) VALUES (?,?,?,?,?)",customer.CustomerId,
         customer.Name,customer.Email,customer.Phone,customer.AccountType)
	if err!=nil{
	    return customerrors.NewRepoError("SaveCustomer",err)
	}
	return nil
}



type AccountSqlDB struct {
	db *sql.DB
}

func NewAccount(db *sql.DB)*AccountSqlDB{
   return &AccountSqlDB{db:db}
}

func(d *AccountSqlDB) GetPin(accountNo string)(string,error){
	var pin string
	log.Printf("account no. is %v",accountNo)
	accountNo1:=strings.TrimSpace(accountNo)
	err:=d.db.QueryRow("SELECT Pin from Account WHERE AccountNo=?",accountNo1).Scan(&pin)
	if err!=nil{
		log.Printf("Invalid User %v",err)
		return "",customerrors.NewRepoError("Invalid User",err)
	}
	log.Println(pin)
	return pin,nil
}

func(d *AccountSqlDB)ChangePin(accountNo string,Pin string)error{

	result,err:=d.db.Exec("UPDATE Account SET Pin=? WHERE AccountNo=?",Pin,accountNo)
	if err!=nil{
		return err
	}
	rowsAffected,err:=result.RowsAffected()
    if err!=nil{
		log.Println(err)
		return err
	}

	if rowsAffected==0{
		return fmt.Errorf("no acccount with this accpuntNo %v",accountNo)
	}

	return nil
}

func(d *AccountSqlDB)SaveAccount(AccountNo string,customerId string,accountType string,Balance float64,Pin string)(error){
    
	_ ,err:=d.db.Exec("INSERT INTO Account(AccountNo,CustomerId,AccountType,Balance,Pin) VALUES ( ? , ? , ? , ? , ? )",AccountNo,customerId,
	accountType,Balance,Pin)

	if err!=nil{
		 log.Printf("failed to create an account %v",err)
		 return customerrors.NewRepoError("SaveAccount",err)
	}

	return nil
}

func (d *AccountSqlDB)GetAccountDetails(AccountNo string)(domain.Account,error){

	var Account domain.Account
	rows:=d.db.QueryRow("SELECT * From Account WHERE AccountNo=?",AccountNo)
	err:=rows.Scan(&Account.AccountNo,&Account.CustomerId,&Account.AccountType,&Account.Balance,&Account.Pin)
	if err!=nil{
		return domain.Account{},customerrors.NewRepoError("GetAccountDetais",err)
	}
	return Account,nil
}

func(d *AccountSqlDB)GetBalance(accountNo string)(float64,error){
	var balance float64
	err:=d.db.QueryRow("SELECT Balance from Account WHERE AccountNo=?",accountNo).Scan(&balance)
	if err!=nil{
		return 0,customerrors.NewRepoError("GetBalance",err)
	}
 
	return balance,nil
}

func(d *AccountSqlDB)SaveBalance(accountNo string,amount float64)error{

	Tx,err:=d.Begin();if err!=nil{
		return customerrors.NewRepoError("SaveAccount",err)
	}
    tx:=Tx.(*SqlDBTransaction).Tx

   result,err:=tx.Exec("UPDATE Account SET Balance=? Where AccountNo=?",amount,accountNo)
	if err!=nil{
		return customerrors.NewRepoError("SaveBalance: query Unsuccessfull",err)
	}
	rowsAffected,err:=result.RowsAffected()
    if err!=nil{
		return customerrors.NewRepoError("SaveBalance",err)
	}

	if rowsAffected==0{
		return customerrors.NewRepoError("SaveBalance",err)
	}

	return nil
}

type TransactionSqlDB struct{
	db *sql.DB
}

func NewTransaction(db *sql.DB)*TransactionSqlDB{
	return &TransactionSqlDB{db:db}
}

func(t *TransactionSqlDB)SaveTransaction(transactionId string,fromAccountNo string,toAccountNo string,Amount float64,time string,status string)error{

	_ ,err:=t.db.Exec("INSERT INTO Transaction(TransactionId,FromAccountId,ToAccountId,AMOUNT,TimeStamp,Status) VALUES (?,?,?,?,?,?)",transactionId,fromAccountNo,
	toAccountNo,Amount,time,status)

	if err!=nil{
		 log.Printf("failed to do transaction %v",err)
		 return customerrors.NewRepoError("saveTransaction: Failed to do Transaction",err)
	}

	return nil
}

func(t *TransactionSqlDB)GetTransactionDetail(transactionId string)(domain.Transaction,error){

	var Transaction domain.Transaction

	rows:=t.db.QueryRow("SELECT * from Transaction WHERE TransactionId=?",transactionId)
	err:=rows.Scan(&Transaction.TransactionId,&Transaction.FromAccountId,&Transaction.ToAccountId,
		&Transaction.Amount,&Transaction.TimeStamp,&Transaction.Status)

	if err!=nil{
		log.Printf("Unable to get pin from the database %v",err)
		return domain.Transaction{},customerrors.NewRepoError("GetTransactionDetail",err)
	}
 
	return Transaction,nil
}



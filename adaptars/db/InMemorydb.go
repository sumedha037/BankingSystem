package db

import (
	"BankingSystem/Core/domain"
	"BankingSystem/Core/ports"
	"fmt"
)


type CustomerDB struct{
	Customer map[string]domain.Customer
}
func NewCustomerDB() *CustomerDB{
	return &CustomerDB{
		Customer:make(map[string]domain.Customer) ,
	    }
}

func (c *CustomerDB)SaveCustomer(customer domain.Customer)error{
	c.Customer[customer.CustomerId]=customer
	return nil
}



type AccountDB struct {
	Account map[string]domain.Account
}
func NewAccountDB()*AccountDB{
	return &AccountDB{
		make(map[string]domain.Account),
	}
}

func(d *AccountDB) GetPin(accountNo string)(string,error){
    _,exist:=d.Account[accountNo];if !exist{
      return "",fmt.Errorf("account no.does not exist %v",accountNo)
	}
	return d.Account[accountNo].Pin,nil
}

func(d *AccountDB)ChangePin(accountNo string,Pin string)error{
	  account,exist:=d.Account[accountNo];if !exist{
      return fmt.Errorf("account no.does not exist %v",accountNo)
	}
	account.Pin=Pin
	d.Account[accountNo]=account
	return nil
}

func(d *AccountDB)CreateAccount(AccountNo string,customerId string,accountType string,Balance float64,Pin string)(domain.Account,error){

	var account = domain.Account{
      AccountNo: AccountNo,
	  CustomerId: customerId,
	  AccountType: accountType,
	  Balance: Balance,
	  Pin: Pin,
	}
    
	d.Account[AccountNo]=account

	return account,nil
}

func(d *AccountDB)GetBalance(accountNo string)(float64,error){

	 _,exist:=d.Account[accountNo];if !exist{
      return 0,fmt.Errorf("account no.does not exist %v",accountNo)
	}
	return d.Account[accountNo].Balance,nil
}

func(d *AccountDB)SaveBalance(accountNo string,amount float64)error{

	account,exist:=d.Account[accountNo];if !exist{
      return fmt.Errorf("account no.does not exist %v",accountNo)
	}
	account.Balance=amount
	d.Account[accountNo]=account
	return nil
}

func (d *AccountDB)GetAccountDetails(accountNo string)(domain.Account,error){

	account,exist:=d.Account[accountNo];if !exist{
		return domain.Account{},fmt.Errorf("account does not exist %v",accountNo)
	}
	return account,nil
}

func(d *AccountDB)SaveAccount(AccountNo string,customerId string,accountType string,Balance float64,Pin string)error{
	var account = domain.Account{
		AccountNo:AccountNo,
		CustomerId: customerId,
		AccountType: accountType,
		Balance: Balance,
		Pin: Pin,
	}
	d.Account[AccountNo]=account
	return nil
}


func(d *AccountDB)Begin()(ports.Transaction,error){
   return &SqlDBTransactionDB{},nil
}

type SqlDBTransactionDB struct{}


func (s *SqlDBTransactionDB)Commit()error{
	return nil
}

func (s *SqlDBTransactionDB)Rollback()error{
	return nil
}


type TransactionDB struct{
	Transaction map[string]domain.Transaction
}
func NewTransactionDB()*TransactionDB{
	return &TransactionDB{
		make(map[string]domain.Transaction),
	}
}

func(t *TransactionDB)SaveTransaction(transactionId string,fromAccountNo string,toAcountNo string,Amount float64,time string,status string)error{
	
	var Transaction = domain.Transaction{
      TransactionId: transactionId,
	  FromAccountId: fromAccountNo,
	  ToAccountId: toAcountNo,
	  Amount: Amount,
	  TimeStamp: time,
	  Status: status,
	}

	t.Transaction[transactionId]=Transaction

	return nil
}

func(t *TransactionDB)GetTransactionDetail(transactionId string)(domain.Transaction,error){

	transaction,exist:=t.Transaction[transactionId];if !exist{
      return domain.Transaction{},fmt.Errorf("transaction Id does not exist %v",transactionId)
	}

	return transaction,nil
}


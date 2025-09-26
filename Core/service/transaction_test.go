package service

import (
	"BankingSystem/Core/domain"
	adaptars "BankingSystem/adaptars/db"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestWithdraw(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	
	b:=NewTransactionService(AccountDB,nil,Helper,nil)

	var acccount = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}

	AccountDB.Account["ab123"]=acccount

	err:=b.Withdraw("ab123",300,"000123")

	if err!=nil{
		t.Errorf("expected no error but got %v",err)
	}

}


func TestDeposite(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	TransactionDB:=adaptars.NewTransactionDB()
	Helper:=NewHelperService(AccountDB)
	
	b:=NewTransactionService(AccountDB,TransactionDB,Helper,nil)

	var acccount = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}

	AccountDB.Account["ab123"]=acccount

	err:=b.Deposite("ab123",300,"000123")

	if err!=nil{
		t.Errorf("expected no error but got %v",err)
	}

}


func TestBalance(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	
	b:=NewBankingService(AccountDB,nil,nil,Helper,nil)

	var acccount = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}

	AccountDB.Account["ab123"]=acccount

	val,err:=b.Balance("ab123","000123")

	if err!=nil{
		t.Errorf("expected no error but got %v",err)
	}

	if val!=30000{
		t.Errorf("expected %v got %v",30000,val)
	}

}



func TestSetPin(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	
	b:=NewBankingService(AccountDB,nil,nil,Helper,nil)

	var acccount = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}

	AccountDB.Account["ab123"]=acccount

	err:=b.SetPin("ab123","000123","123456")

	if err!=nil{
		t.Errorf("expected no error but got %v",err)
	}

	if AccountDB.Account["ab123"].Pin!="123456"{
		t.Errorf("expected %v",123456)
	}
}




func TestTransfer(t *testing.T){
	TransactionDB:=adaptars.NewTransactionDB()
	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	IdGenerator:=NewIdGenerator()
	
	b:=NewTransactionService(AccountDB,TransactionDB,Helper,IdGenerator)

	var account1 = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}

	var account2 = domain.Account{
       AccountNo: "bc234",
	   CustomerId: "cust2",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000162",
	}

   AccountDB.Account["ab123"]=account1
   AccountDB.Account["bc234"]=account2

   s,err:=b.Transfer("ab123","000123","bc234",4000)

   if err!=nil{
      t.Errorf("expected no error but got %v",err)
   }
   if s==""{
       t.Error("expected randomstring of length 8 but got an empty string")
   }

}


func TestCreateAccount(t *testing.T){
	CustomerDB:=adaptars.NewCustomerDB()
	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	IdGenerator:=NewIdGenerator()
	TransactionDB:=adaptars.NewTransactionDB()
	
	b:=NewBankingService(AccountDB,CustomerDB,TransactionDB,Helper,IdGenerator)

	var customer= domain.Customer{
		CustomerId: "cust1",
		Name: "Shreya",
		Email: "abc@gamil.com",
		Phone: "887687574",
		AccountType: "Saving Account",
	}

	
	CustomerDB.Customer["cust1"]=customer

	account:=b.CreateAccount(customer)
	if account.CustomerId!="cust1"{
		t.Error("expected cust1")
	}

	AccountDB.Account[account.AccountNo]=account
}



func TestValidateUser(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)
	
	b:=NewBankingService(AccountDB,nil,nil,Helper,nil)


	var account1 = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}
	
	AccountDB.Account["ab123"]=account1

	_,err:=b.HelperService.ValidateUser("ab123","000123")
	assert.NoError(t,err,"Expected true")

	_,err=b.HelperService.ValidateUser("ab123","122456")
     assert.Error(t,err,"Expected error")
}




func TestIncreaseAmount(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)

	b:=NewBankingService(AccountDB,nil,nil,Helper,nil)

	var account1 = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}
    AccountDB.Account["ab123"]=account1

	err:=b.HelperService.IncreaseAmount("ab123",1000)
	if err!=nil{
		t.Error("error not expected")
	}
	if AccountDB.Account["ab123"].Balance!=31000{
		t.Errorf("expected non equal balance")
	}
}




func TestDecreaseAmount(t *testing.T){

	AccountDB:=adaptars.NewAccountDB()
	Helper:=NewHelperService(AccountDB)

	b:=NewBankingService(AccountDB,nil,nil,Helper,nil)

	var account1 = domain.Account{
       AccountNo: "ab123",
	   CustomerId: "cust1",
	   AccountType: "Saving Account",
	   Balance:30000,
	   Pin: "000123",
	}
    AccountDB.Account["ab123"]=account1
	
	err:=b.HelperService.DecreaseAmount("ab123",1000)
	if err!=nil{
		t.Error("error not expected")
	}
	if AccountDB.Account["ab123"].Balance!=29000{
		t.Errorf("expected non equal balance")
	}
	
}
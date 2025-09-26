package service

import (
	"BankingSystem/Core/domain"
	"BankingSystem/Core/ports"
	customerrors "BankingSystem/customErrors"
	"log"
)

type BankingService struct{
     AccountRepo ports.AccountRepository
	 CustomerRepo  ports.CustomerRepository
	 TransactionRepo ports.TransactionRepository
	 HelperService   ports.Helper
	 IdGenerator   ports.IdGenerator
}


func NewBankingService(a ports.AccountRepository,c ports.CustomerRepository,t ports.TransactionRepository,h ports.Helper,i ports.IdGenerator)*BankingService{
	return &BankingService{
		AccountRepo: a,
		CustomerRepo: c,
		TransactionRepo: t,
		HelperService: h,
		IdGenerator: i,
	}
}



func(b *BankingService)SetPin(accountNo string,OldPin string, NewPin string)error{

	ok,err:=b.HelperService.ValidateUser(accountNo,OldPin);if !ok{
		return customerrors.NewServiceError("SetPin:Unauthorized User",err)
	}
  
	err=b.AccountRepo.ChangePin(accountNo,NewPin)
	if err!=nil{
		return customerrors.NewServiceError("Change Pin",err)
	}
	return nil
}


func(b *BankingService)CreateAccount(customer domain.Customer)domain.Account{
	var account domain.Account
   err:= b.CustomerRepo.SaveCustomer(customer)
   if err!=nil{
	log.Printf("failed to save customer in database %v",err)
	return domain.Account{}
   }

	AccountNo:=b.IdGenerator.GenerateSequentialID(12)
	Pin:=b.IdGenerator.GenerateSequentialID(6)
	Balance:=0.00

	err=b.AccountRepo.SaveAccount(AccountNo,customer.CustomerId,customer.AccountType,Balance,Pin)
	if err!=nil{
	    log.Printf("Failed to Create an Account %v",err)
		return domain.Account{}
	}

   account,err=b.AccountRepo.GetAccountDetails(AccountNo)
   if err!=nil{
	log.Printf("Failed to Get Account Details for %v",AccountNo)
	return domain.Account{}
   }
    return account
}


func(b *BankingService) Balance(accountno string,Pin string)(float64,error){

	 ok,err:=b.HelperService.ValidateUser(accountno,Pin);if !ok{
		return 0,customerrors.NewServiceError("Balance Unauthorized User",err)
	 }

	 balance,err:=b.AccountRepo.GetBalance(accountno)
	 if err!=nil{
		log.Println("failed to get the balance")
		return 0,customerrors.NewServiceError("Balance",err)
	 }
	 
	 return balance,nil
}
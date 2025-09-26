package ports

import "BankingSystem/Core/domain"

type TransactionService interface {
	Deposite(accountno string, amount float64, Pin string) error
	Withdraw(accountno string, amount float64, Pin string) error
	Transfer(fromAccountNo string, fromAccountPin string, toAcountNo string, Amount float64) (string, error)
}

type AccountsService interface {
	SetPin(accountNo string,OldPin string, NewPin string)error
	CreateAccount(customer domain.Customer) domain.Account
	Balance(accountno string, Pin string) (float64, error)
}

type Helper interface{
	IncreaseAmount(accountNo string,amount float64)error
	DecreaseAmount(accountNo string,amount float64)error
	ValidateUser(accountNo string,Pin string)(bool,error)
}

type IdGenerator interface{
	GenerateSequentialID(length int) string
}
package ports

import "BankingSystem/Core/domain"

type CustomerRepository interface {
	SaveCustomer(customer domain.Customer)error
}

type AccountRepository interface{
    GetPin(accountNo string)(string,error)
	ChangePin(accountNo string,Pin string)error
  	SaveAccount(AccountNo string,customerId string,accountType string,Balance float64,Pin string)(error)
	GetAccountDetails(AccountNo string)(domain.Account,error)
	GetBalance(accountNo string)(float64,error)
	SaveBalance(accountNo string,amount float64)error
}

type TransactionRepository interface{
	SaveTransaction(transactionId string,fromAccountNo string,toAcountNo string,Amount float64,timestamp string,status string)error
	GetTransactionDetail(transactionId string)(domain.Transaction,error)
}

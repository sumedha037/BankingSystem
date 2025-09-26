package domain

type Customer struct {
	CustomerId string
	Name       string
	Email      string
	Phone      string
	AccountType string
}

type Account struct {
	AccountNo   string
	CustomerId  string
	AccountType string
	Balance     float64
	Pin   		string
}

type Transaction struct {
	TransactionId string
	FromAccountId string
	ToAccountId   string
	Amount        float64
	TimeStamp     string
	Status        string
}



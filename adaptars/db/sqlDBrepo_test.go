package db



import (
    "BankingSystem/Core/domain"
    "github.com/stretchr/testify/assert"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
)

func TestSaveCustomer(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewCustomer(db)

    customer := domain.Customer{
        CustomerId:  "test123",
        Name:        "Test User",
        Email:       "test@example.com",
        Phone:       "1234567890",
        AccountType: "Savings",
    }

    mock.ExpectExec("INSERT INTO Customer").
        WithArgs(customer.CustomerId, customer.Name, customer.Email, customer.Phone, customer.AccountType).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.SaveCustomer(customer)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}

func TestGetPin(t *testing.T){
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewAccount(db)

    account:=domain.Account{
		AccountNo: "abc123",
		Pin: "000123",
	}

     rows := sqlmock.NewRows([]string{"Pin"}).AddRow(account.Pin)

    mock.ExpectQuery("SELECT Pin from Account WHERE AccountNo=?").WithArgs(account.AccountNo).WillReturnRows(rows)

    Pin, err := repo.GetPin(account.AccountNo)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }

    if Pin!=account.Pin{
        t.Errorf("InValid User")
    }

    if err:= mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unfulfilled expectations: %v", err)
    }

}


func TestChangePin(t *testing.T){
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewAccount(db)

    account:=domain.Account{
		AccountNo: "abc123",
		CustomerId: "cust1",
		AccountType: "Saving Account",
		Balance: 30000,
		Pin: "000123",
	}
    Pin:="123456"

     mock.ExpectExec(`(?i)UPDATE\s+Account\s+SET\s+Pin\s*=\s*\?\s+WHERE\s+AccountNo\s*=\s*\?`).WithArgs(Pin,account.AccountNo).WillReturnResult(sqlmock.NewResult(1,1))


    err = repo.ChangePin(account.AccountNo,Pin)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }


    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}


func TestSaveAccount(t *testing.T){
 db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()   

     account:=domain.Account{
		AccountNo: "abc123",
		CustomerId: "cust1",
		AccountType: "Saving Account",
		Balance: 30000,
		Pin: "000123",
	}

    repo := NewAccount(db)

    mock.ExpectExec("INSERT INTO Account").
        WithArgs(account.AccountNo,account.CustomerId,account.AccountType,account.Balance,account.Pin).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.SaveAccount(account.AccountNo,account.CustomerId,account.AccountType,account.Balance,account.Pin)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}


func TestGetAccountDetails(t *testing.T){
  db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()    
    
    account:=domain.Account{
		AccountNo: "abc123",
		CustomerId: "cust1",
		AccountType: "Saving Account",
		Balance: 30000,
		Pin: "000123",
	}

mock.ExpectQuery(`(?i)SELECT\s+\*\s+FROM\s+Account\s+WHERE\s+AccountNo\s*=\s*\?`).
    WithArgs(account.AccountNo).
    WillReturnRows(sqlmock.NewRows([]string{"AccountNo", "CustomerId", "AccountType", "Balance", "Pin"}).
        AddRow(account.AccountNo, account.CustomerId, account.AccountType, account.Balance, account.Pin))



    repo:=NewAccount(db)
    repo.GetAccountDetails(account.AccountNo)

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}



func TestSaveBalance(t *testing.T) {
	
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE Account SET Balance=\\? Where AccountNo=\\?").
		WithArgs(5000.0, "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	accDB := &AccountSqlDB{db: db}


	err = accDB.SaveBalance("abc123", 5000.0)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransactionDetails(t *testing.T){
      db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()    
   	  Transaction := domain.Transaction{
      TransactionId: "abc123",
	  FromAccountId: "00000000000012",
	  ToAccountId:"000000000011",
	  Amount: 2000,
	  TimeStamp: "2006-01-02 15:04:05",
	  Status: "Successfull",
	}

mock.ExpectQuery(`(?i)SELECT\s+\*\s+FROM\s+Transaction\s+WHERE\s+TransactionId\s*=\s*\?`).
    WithArgs().
    WillReturnRows(sqlmock.NewRows([]string{"TransactionId", "FromAccountId", "ToAccountId", "Amount", "TimeStamp","Status"}).
        AddRow(Transaction.TransactionId, Transaction.FromAccountId, Transaction.ToAccountId, Transaction.Amount, Transaction.TimeStamp,Transaction.Status))

    repo:=NewTransaction(db)
    repo.GetTransactionDetail(Transaction.TransactionId)

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}


func TestSaveTransaction(t *testing.T){
      db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()    
   	  Transaction := domain.Transaction{
      TransactionId: "abc123",
	  FromAccountId: "00000000000012",
	  ToAccountId:"000000000011",
	  Amount: 2000,
	  TimeStamp: "2006-01-02 15:04:05",
	  Status: "Successfull",
	}

     mock.ExpectExec("INSERT INTO Transaction").
        WithArgs(Transaction.TransactionId, Transaction.FromAccountId, Transaction.ToAccountId, Transaction.Amount, Transaction.TimeStamp,Transaction.Status).
        WillReturnResult(sqlmock.NewResult(1, 1))

         repo:=NewTransaction(db)
     repo.SaveTransaction(Transaction.TransactionId, Transaction.FromAccountId, Transaction.ToAccountId, Transaction.Amount, Transaction.TimeStamp,Transaction.Status)

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }    


}


func TestGetBalance(t *testing.T){
      db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create sqlmock: %v", err)
    }
    defer db.Close()

    repo := NewAccount(db)

    account:=domain.Account{
		AccountNo: "abc123",
        Balance:30000.00,
	}

     rows := sqlmock.NewRows([]string{"Balance"}).AddRow(account.Balance)

    mock.ExpectQuery(`(?i)SELECT\s+Balance\s+from\s+Account\s+WHERE\s+AccountNo\s*=\s*\?`).WithArgs(account.AccountNo).WillReturnRows(rows)
    balance, err := repo.GetBalance(account.AccountNo)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }

    if balance != account.Balance {
        t.Errorf("Expected balance %v, got %v", account.Balance, balance)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Unmet expectations: %v", err)
    }
}



package handlers

 import (
	"BankingSystem/Core/domain"
	"BankingSystem/middleware"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/mock"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
 )




type mockTransactionService struct{
	mock.Mock
}

func (m *mockTransactionService)Deposite(accountno string, amount float64, Pin string) error{
	args := m.Called(accountno,amount,Pin)
	return args.Error(0)
}

func (m *mockTransactionService)Withdraw(accountno string, amount float64, Pin string) error{
	args := m.Called(accountno,amount,Pin)
	return args.Error(0)
}

func (m *mockTransactionService)Transfer(fromAccountNo string, fromAccountPin string, toAcountNo string, Amount float64) (string, error){
	args := m.Called(fromAccountNo,fromAccountPin,toAcountNo,Amount)
	return args.Get(0).(string), args.Error(1)
}


type mockAccountService struct{
	mock.Mock
}

func(m *mockAccountService)SetPin(accountNo string,OldPin string, NewPin string)error{
	args := m.Called(accountNo,OldPin,NewPin)
	return args.Error(0)
}

func(m *mockAccountService)CreateAccount(customer domain.Customer) domain.Account{
	args:=m.Called(customer)
	return args.Get(0).(domain.Account)
}

func(m *mockAccountService)Balance(accountno string, Pin string) (float64, error){
	args:=m.Called(accountno,Pin)
	return args.Get(0).(float64),args.Error(1)
}

type mockHelperService struct{
	mock.Mock
}

func (m *mockHelperService)IncreaseAmount(accountNo string,amount float64)error{
	args := m.Called(accountNo,amount)
	return args.Error(0)
}

func (m *mockHelperService)DecreaseAmount(accountNo string,amount float64)error{
	args := m.Called(accountNo,amount)
	return args.Error(0)
}

func (m *mockHelperService)ValidateUser(accountNo string,Pin string)(bool,error){
	args := m.Called(accountNo,Pin)
	return args.Get(0).(bool),args.Error(1)
}




func setup() (*mockTransactionService, *mockAccountService, *mockHelperService,*Handlers) {
	mockTxn := new(mockTransactionService)
	mockAcc := new(mockAccountService)
	mockHelper := new(mockHelperService)
	h := NewHandler(mockTxn, mockAcc, mockHelper)
	return mockTxn, mockAcc, mockHelper,h
}

func TestCheckBalance(t *testing.T) {
	_, mockAcc, _,h:= setup()


	mockAcc.On("Balance", "abc123", "000123").Return(5000.0, nil)

	body, _ := json.Marshal(map[string]string{"Pin": "000123"})
	req := httptest.NewRequest(http.MethodPost, "/CheckBalance", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.AccountKey, "abc123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.CheckBalance(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	
	mockAcc.On("Balance", "abc123", "000124").Return(0.0, errors.New("invalid pin"))

	body2, _ := json.Marshal(map[string]string{"Pin": "000124"})
	req2 := httptest.NewRequest(http.MethodPost, "/CheckBalance", bytes.NewReader(body2))
	ctx2 := context.WithValue(req2.Context(), middleware.AccountKey, "abc123")
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	h.CheckBalance(w2, req2)

	resp2 := w2.Result()
	assert.Equal(t, http.StatusInternalServerError, resp2.StatusCode)

	mockAcc.AssertExpectations(t)
}

func TestWithdrawAmount(t *testing.T) {
	mockTxn, mockAcc, _, h := setup()

	mockTxn.On("Withdraw", "abc123", 1000.0, "000123").Return(nil)

	body, _ := json.Marshal(map[string]interface{}{"Pin": "000123", "Amount": 1000})
	req := httptest.NewRequest(http.MethodPost, "/Withdraw", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.AccountKey, "abc123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.WithdrawAmount(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockTxn.On("Withdraw", "abc123", 2000.0, "000999").Return(errors.New("insufficient balance"))

	body2, _ := json.Marshal(map[string]interface{}{"Pin": "000999", "Amount": 2000})
	req2 := httptest.NewRequest(http.MethodPost, "/Withdraw", bytes.NewReader(body2))
	ctx2 := context.WithValue(req2.Context(), middleware.AccountKey, "abc123")
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	h.WithdrawAmount(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Result().StatusCode)

	mockTxn.AssertExpectations(t)
	mockAcc.AssertExpectations(t)
}

func TestDepositeAmount(t *testing.T) {
	mockTxn, mockAcc, _, h := setup()

	mockTxn.On("Deposite", "abc123", 1000.0, "000123").Return(nil)

	body, _ := json.Marshal(map[string]interface{}{"Pin": "000123", "Amount": 1000})
	req := httptest.NewRequest(http.MethodPost, "/Deposite", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.AccountKey, "abc123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.DepositeAmount(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockTxn.On("Deposite", "abc123", 1000.0, "000999").Return(errors.New("invalid pin"))

	body2, _ := json.Marshal(map[string]interface{}{"Pin": "000999", "Amount": 1000})
	req2 := httptest.NewRequest(http.MethodPost, "/Deposite", bytes.NewReader(body2))
	ctx2 := context.WithValue(req2.Context(), middleware.AccountKey, "abc123")
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	h.DepositeAmount(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Result().StatusCode)

	mockTxn.AssertExpectations(t)
	mockAcc.AssertExpectations(t)
}

func TestTransferAmount(t *testing.T) {
	mockTxn, _, _, h := setup()

	mockTxn.On("Transfer", "abc123", "000123", "abc124", 3000.0).Return("Transfer Success", nil)

	body, _ := json.Marshal(map[string]interface{}{"FromAccountPin": "000123", "ToAccountNo": "abc124", "Amount": 3000})
	req := httptest.NewRequest(http.MethodPost, "/Transfer", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.AccountKey, "abc123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.TransferAmount(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	mockTxn.On("Transfer", "abc123", "000124", "abc124", 3000.0).Return("", errors.New("wrong pin"))

	body2, _ := json.Marshal(map[string]interface{}{"FromAccountPin": "000124", "ToAccountNo": "abc124", "Amount": 3000})
	req2 := httptest.NewRequest(http.MethodPost, "/Transfer", bytes.NewReader(body2))
	ctx2 := context.WithValue(req2.Context(), middleware.AccountKey, "abc123")
	req2 = req2.WithContext(ctx2)

	w2 := httptest.NewRecorder()
	h.TransferAmount(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Result().StatusCode)

	mockTxn.AssertExpectations(t)
}

func TestSetPin(t *testing.T) {
	_, mockAcc, _, h := setup()


	mockAcc.On("SetPin", "abc123", "000123", "123456").Return(nil)

	body, _ := json.Marshal(map[string]string{"AccountNo": "abc123", "OldPin": "000123", "NewPin": "123456"})
	req := httptest.NewRequest(http.MethodPost, "/SetPin", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.SetPin(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)


	mockAcc.On("SetPin", "abc123", "111111", "123456").Return(errors.New("wrong old pin"))

	body2, _ := json.Marshal(map[string]string{"AccountNo": "abc123", "OldPin": "111111", "NewPin": "123456"})
	req2 := httptest.NewRequest(http.MethodPost, "/SetPin", bytes.NewReader(body2))
	w2 := httptest.NewRecorder()
	h.SetPin(w2, req2)
	assert.Equal(t, http.StatusInternalServerError, w2.Result().StatusCode)

	mockAcc.AssertExpectations(t)
}

func TestCreateAccount(t *testing.T) {
	_, mockAcc, _, h := setup()

	
	acc := domain.Account{AccountNo: "abc123", CustomerId: "cust1"}
	cust := domain.Customer{CustomerId: "cust1", Name: "Rob"}
	mockAcc.On("CreateAccount", cust).Return(acc)

	body, _ := json.Marshal(cust)
	req := httptest.NewRequest(http.MethodPost, "/CreateAccount", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateAccount(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	req2 := httptest.NewRequest(http.MethodPost, "/CreateAccount", bytes.NewReader([]byte("{invalid_json")))
	w2 := httptest.NewRecorder()
	h.CreateAccount(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Result().StatusCode)

	mockAcc.AssertExpectations(t)
}
package app

import (
	"BankingSystem/Core/service"
	adaptars "BankingSystem/adaptars/db"
	"BankingSystem/adaptars/handlers"
	"BankingSystem/middleware"
	"BankingSystem/dbInstance"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start(){

    database := dbInstance.GetInstance()
 
   CustomerRepo:=adaptars.NewCustomer(database)
   AccountRepo:=adaptars.NewAccount(database)
   TransactionRepo:=adaptars.NewTransaction(database)

   HelperService:=service.NewHelperService(AccountRepo)
   IdGenerator:=service.NewIdGenerator()

   TransactionService:=service.NewTransactionService(AccountRepo,TransactionRepo,HelperService,IdGenerator)
   BankingService:=service.NewBankingService(AccountRepo,CustomerRepo,TransactionRepo,HelperService,IdGenerator)
   
    h:=handlers.NewHandler(TransactionService,BankingService,HelperService)
    r := mux.NewRouter()
    r.HandleFunc("/CreateAccount", h.CreateAccount).Methods(http.MethodPost)
	r.HandleFunc("/SetPin", h.SetPin).Methods(http.MethodPost)
	r.HandleFunc("/Login", h.Login).Methods(http.MethodPost)

	
	r.Handle("/Deposit", middleware.AuthMiddleware(http.HandlerFunc(h.DepositeAmount))).Methods(http.MethodPost)
	r.Handle("/Withdraw", middleware.AuthMiddleware(http.HandlerFunc(h.WithdrawAmount))).Methods(http.MethodPost)
	r.Handle("/Transfer",  middleware.AuthMiddleware(http.HandlerFunc(h.TransferAmount))).Methods(http.MethodPost)
	r.Handle("/CheckBalance",  middleware.AuthMiddleware(http.HandlerFunc(h.CheckBalance))).Methods(http.MethodPost)
 
    log.Println("Server running on:8080")
    http.ListenAndServe(":8080", r)
}
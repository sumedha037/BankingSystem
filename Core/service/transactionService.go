package service

import "BankingSystem/Core/ports"
import "BankingSystem/customErrors"
import "log"
import "fmt"
import "time"

type TransactionService struct {
	Tx ports.TransactionManager
	TransactionRepo  ports.TransactionRepository
	HelperService ports.Helper
	IdGenerator ports.IdGenerator
}


func NewTransactionService(tx ports.TransactionManager,TransactionRepo ports.TransactionRepository,HelperService ports.Helper,IdGenerator ports.IdGenerator)*TransactionService{
     return &TransactionService{
		Tx:tx,
		TransactionRepo: TransactionRepo,
		HelperService: HelperService,
		IdGenerator: IdGenerator,
	 }
}


func(b *TransactionService) Withdraw(accountno string,amount float64,Pin string)error{

	ok,err:=b.HelperService.ValidateUser(accountno,Pin);if !ok {
		return customerrors.NewServiceError("Withdraw: Unauthorized User",err)
	}

	tx,err:=b.Tx.Begin();if err!=nil{
		return customerrors.NewServiceError("Transfer:",err)
	}
   
	defer func(){
		r:=recover();if r!=nil{
			tx.Rollback()
			log.Println(r)
		}
	}()
	    if amount<=0{
			return customerrors.NewServiceError("WithDraw DecreaseAmount",fmt.Errorf("negative amount"))
		}
		err=b.HelperService.DecreaseAmount(accountno,amount)
		if err!=nil{
			return customerrors.NewServiceError("WithDraw DecreaseAmount",err)
		}
		tx.Commit()

	 return nil
}

func(b *TransactionService) Deposite(accountno string,amount float64,Pin string)error{
	ok,err:=b.HelperService.ValidateUser(accountno,Pin);if!ok{
		return customerrors.NewServiceError("Withdraw:Unauthorized User",err)
	}

	tx,err:=b.Tx.Begin();if err!=nil{
		return customerrors.NewServiceError("Transfer:",err)
	}
   
	defer func(){
		r:=recover();if r!=nil{
			tx.Rollback()
			log.Println(r)
		}
	}()

	if amount<=0{
		return customerrors.NewServiceError("Deposite Incraese Amount",fmt.Errorf("negative amount"))
	}

	err=b.HelperService.IncreaseAmount(accountno,amount)
	    if err!=nil{
			return customerrors.NewServiceError("Deposite Incraese Amount",err)
		}
	tx.Commit()
  return nil
}

func(b *TransactionService) Transfer(fromAccountNo string,fromAccountPin string,toAcountNo string,Amount float64)(string,error){
    var status string

	if fromAccountNo==toAcountNo{
		return "",customerrors.NewServiceError("Transfer: r",fmt.Errorf("cannot transfer money in same account"))
	}

	ok,err:=b.HelperService.ValidateUser(fromAccountNo,fromAccountPin);if !ok{
		return "",customerrors.NewServiceError("Transfer:Unauthorized User",err)
	}


	tx,err:=b.Tx.Begin();if err!=nil{
		return "",customerrors.NewServiceError("Transfer:",err)
	}
	defer tx.Rollback()

	defer func(){
		r:=recover();if r!=nil{
			tx.Rollback()
			log.Println(r)
		}
	}()

	id:=b.IdGenerator.GenerateSequentialID(8)

	timestamp:=time.Now()
	formattedTime := timestamp.Format("2006-01-02 15:04:05")

	err=b.HelperService.DecreaseAmount(fromAccountNo,Amount)
	if err!=nil{
		tx.Rollback()
		return "",customerrors.NewServiceError("transfer",err)
	}
	err=b.HelperService.IncreaseAmount(toAcountNo,Amount)
	if err!=nil{
		tx.Rollback()
		return "",customerrors.NewServiceError("transfer",err)
	}

	status="Successfull"

	tx.Commit()

	b.TransactionRepo.SaveTransaction(id,fromAccountNo,toAcountNo,Amount,formattedTime,status)
	 return id,nil
}
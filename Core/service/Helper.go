package service

import (
	customerrors "BankingSystem/customErrors"
  "BankingSystem/Core/ports"
	"fmt"
	"log"
	"os"
	"strconv"
)


type HelperService struct{
  AccountRepo ports.AccountRepository
}

func NewHelperService(AccountRepo ports.AccountRepository)*HelperService{
  return &HelperService{AccountRepo:AccountRepo}
}

type IdGenerator struct{}

func NewIdGenerator()*IdGenerator{
    return &IdGenerator{}
}



func(b *HelperService) ValidateUser(accountNo string,Pin string)(bool,error){

   s,err:=b.AccountRepo.GetPin(accountNo)
   if err!=nil{
    return false,customerrors.NewServiceError("ValidateUser",err)
   }

   if s!=Pin{
    return false,customerrors.NewServiceError("ValidateUser",fmt.Errorf("unauthorized User"))
   }
   log.Println(s)
 return true,nil
}



func (b *HelperService) IncreaseAmount(accountNo string,amount float64)error{
   currentAmount,err:=b.AccountRepo.GetBalance(accountNo)
  if err!=nil{
    log.Printf("unable to get current balance %v",err)
    return customerrors.NewServiceError("Increase Amount",err)
  }

  if amount<0{
    return customerrors.NewServiceError("IncreaseAmount",fmt.Errorf("amount less than zero"))
  }
  currentAmount+=amount
  return b.AccountRepo.SaveBalance(accountNo,currentAmount)
}



func (b *HelperService) DecreaseAmount(accountNo string,amount float64)error{
   currentAmount,err:=b.AccountRepo.GetBalance(accountNo)
  if err!=nil{
    log.Printf("unable to get current balance %v",err)
    return customerrors.NewServiceError("DecreaseAmount",err)
  }
  if currentAmount<0{
    return customerrors.NewServiceError("DecreaseAmount",fmt.Errorf("amount less than zero"))
  }
  if currentAmount<amount{
    return customerrors.NewServiceError("DecreaseAmount",fmt.Errorf("insufficient Amount"))
  }
  currentAmount-=amount
  return b.AccountRepo.SaveBalance(accountNo,currentAmount)
}

func (b *IdGenerator) GenerateSequentialID(length int) string {

	counter := 0
	data, err := os.ReadFile("counter.txt")
	if err == nil {
		val, convErr := strconv.Atoi(string(data))
		if convErr == nil {
			counter = val
		}
	}


	counter++

	_ = os.WriteFile("counter.txt", []byte(strconv.Itoa(counter)), 0644)

	num := strconv.Itoa(counter)
	for len(num) < length {
		num = "0" + num
	}
	return num
}

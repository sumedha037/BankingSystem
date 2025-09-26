package handlers

import (
	"BankingSystem/Core/service"
	"encoding/json"
	"net/http"
	"log"
)

func (h *Handlers) Login(w http.ResponseWriter,r *http.Request){
   var input struct{
	AccountNo   string
	Pin         string
   }

    if err:=json.NewDecoder(r.Body).Decode(&input);err!=nil{
		log.Println("Decode error:", err)
         http.Error(w,"Failed to Decode data",http.StatusBadRequest)
		 return
	}

	ok,err:=h.helperService.ValidateUser(input.AccountNo,input.Pin);if !ok {
		 http.Error(w,err.Error(),http.StatusUnauthorized)
	   return
	}

   	Token,err:=service.GenerateJWT(input.AccountNo)
	if err!=nil{
       http.Error(w,"Failed to generate Token",http.StatusUnauthorized)
	   return
	}

	type Response struct{
		Token string
	}

	resp:=Response{
		Token: Token,
	}

	w.Header().Set("content-type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

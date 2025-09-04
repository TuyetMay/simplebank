package api

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

type transferRequest struct{
	FromAccountID    int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID      int64 `json:"to_account_id" binding:"required,min=1"`
	Amount           int64 `json:"amount" binding:"required,gt=0"`
	Currency         string `json:"currency" binding:"required,currency"`
}

func(server *Server) createTransfer(ctx *gin.Context){
	var req transferRequest
	
	log.Println("Transfer request received")
	
	if err := ctx.ShouldBindJSON(&req); err != nil{
		log.Printf("Binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	log.Printf("Request parsed: %+v", req)

	// check the validity of the request.fromAccountID and currency
	if !server.validAccount(ctx, req.FromAccountID, req.Currency){
		log.Printf("Invalid from account: %d", req.FromAccountID)
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency){
		log.Printf("Invalid to account: %d", req.ToAccountID)
		return
	}

	log.Println("Both accounts validated successfully")

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	log.Printf("Calling TransferTx with params: %+v", arg)

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil{
		log.Printf("TransferTx error: %v", err)
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	log.Println("Transfer completed successfully")
	ctx.JSON(http.StatusOK, result)
}

func(server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool{
	log.Printf("Validating account %d with currency %s", accountID, currency)
	
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil{
		log.Printf("GetAccount error for ID %d: %v", accountID, err)
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	log.Printf("Account found: %+v", account)

	if account.Currency != currency{
		err := fmt.Errorf("account[%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		log.Printf("Currency mismatch: %v", err)
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}
	
	return true
}
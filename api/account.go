package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"

)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency"binding:"required, oneof=USD EUR"`
}

func(server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	} // ddojc JSON body -> map vào struct -> validate theo tag -> trả lỗi nếu có


	arg := db.CreateAccountParams{
		Owner: req.Owner, // lấy từ json client gửi lên
		Currency: req.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx,arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}
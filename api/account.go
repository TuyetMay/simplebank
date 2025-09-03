package api

import (
	"net/http"
    "database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"

)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func(server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err !=nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
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

type getAccountRequest struct {
	ID int64 `uri:"id" binding: "required,min=1"` // id này là trường bắt buộc, số nhỏ nhất là 1
}

func (server *Server) getAccount(ctx *gin.Context){ //gin.Context chứ tất cả thông tin về request,response,...
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON (http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required, min=5,max=10"`
}

func(server *Server) ListAccount(ctx *gin.Context){
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize, // bỏ qua mấy trang đầu 
	}

	accounts, err := server.store.ListAccounts(ctx.Request.Context(),arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}
package api

import (
	"net/http"
    "database/sql"
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/lib/pq"

)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
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
		if pqErr, ok := err.(*pq.Error); ok{
			switch pqErr.Code.Name(){
			case "foreign_key_violation","unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type getAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"` // Correct
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
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"` // Fixed: removed space after required
}

func(server *Server) ListAccount(ctx *gin.Context){
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return // Fixed: added missing return
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
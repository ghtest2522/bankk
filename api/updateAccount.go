package api

import (
	db "bank/db/sqlc"
	"bank/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateAccountRequest struct {
	ID     int64 `json:"id" binding:"numeric,required"`
	Amount int64 `json:"amount" binding:"numeric,required"`
}

func (server *HttpServer) updateAccountBalance(ctx *gin.Context) {
	var req UpdateAccountRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		SendError(ctx, http.StatusBadRequest, err)
		return
	}

	args := db.AddBalanceToAccountParams{
		ID:     req.ID,
		Amount: req.Amount,
	}
	result, err := server.store.AddBalanceToAccount(ctx, args)

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	SendOKRespnse(ctx, util.AccountWasUpdated, result)
}

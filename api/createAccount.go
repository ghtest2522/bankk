package api

import (
	db "bank/db/sqlc"
	"bank/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *HttpServer) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	err := ctx.BindJSON(&req)

	if err != nil {
		server.resultSender.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	result, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		server.resultSender.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	server.resultSender.SendOKRespnse(ctx, util.AccountWasCeated, result)
}

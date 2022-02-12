package api

import (
	"bank/util"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,numeric,min=1"`
}

func (server *HttpServer) getAccount(ctx *gin.Context) {
	var req GetAccountRequest
	err := ctx.BindUri(&req)

	if err != nil {
		SendError[any](ctx, http.StatusBadRequest, err)
		return
	}

	result, err := server.store.GetAccount(ctx, req.ID)
	if err == sql.ErrNoRows {
		SendError[any](ctx, http.StatusNotFound, err)
		return
	}

	if err != nil {
		SendError[any](ctx, http.StatusInternalServerError, err)
		return
	}

	SendOKRespnse(ctx, util.AccountWasFound, result)
}

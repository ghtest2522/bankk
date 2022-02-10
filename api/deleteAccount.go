package api

import (
	"bank/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *HttpServer) deleteAccount(ctx *gin.Context) {
	var req DeleteAccountRequest
	err := ctx.BindUri(&req)
	if err != nil {
		server.resultSender.SendError(ctx, http.StatusBadRequest, err)
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		server.resultSender.SendError(ctx, http.StatusInternalServerError, err)
		return
	}

	server.resultSender.SendOKRespnse(ctx, util.AccountWasDeleted, nil)
}

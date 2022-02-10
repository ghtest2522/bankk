package api

import (
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	result, err := server.store.GetAccount(ctx, req.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, result)

}

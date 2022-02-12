package api

import (
	db "bank/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenericType interface{ db.Account | db.Entry }

type ResponseResult[T GenericType] struct {
	Message string
	Data    T
}

type ResponseErrorResult struct {
	Message string
}

func createError(err error) ResponseErrorResult {
	return ResponseErrorResult{Message: err.Error()}
}

func SendError(ctx *gin.Context, httpStatus int, err error) {
	errorMsg := createError(err)
	ctx.JSON(httpStatus, errorMsg)
}

func createOKRespone[T GenericType](msg string, data T) ResponseResult[T] {
	return ResponseResult[T]{Message: msg, Data: data}
}

func SendOKRespnse[T GenericType](ctx *gin.Context, msg string, data T) {
	response := createOKRespone(msg, data)

	ctx.JSON(http.StatusOK, response)
}

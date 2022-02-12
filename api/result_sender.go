package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseResult[T any] struct {
	Message string
	Data    T
}

func createError[T any](err error) ResponseResult[T] {
	return ResponseResult[T]{Message: err.Error()}
}

func SendError[T any](ctx *gin.Context, httpStatus int, err error) {
	errorMsg := createError[T](err)
	ctx.JSON(httpStatus, errorMsg)
}

func createOKRespone[T any](msg string, data T) ResponseResult[T] {
	return ResponseResult[T]{Message: msg, Data: data}
}

func SendOKRespnse[T any](ctx *gin.Context, msg string, data T) {
	response := createOKRespone(msg, data)

	ctx.JSON(http.StatusOK, response)
}

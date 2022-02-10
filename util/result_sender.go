package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseResult struct {
	Message string
	Data    interface{}
}

type ResponseSender interface {
	createError(err error) ResponseResult
	SendError(ctx *gin.Context, httpStatus int, err error)
	createOKRespone(msg string, data interface{}) ResponseResult
	SendOKRespnse(ctx *gin.Context, msg string, data interface{})
}

type JSONResponseSender struct{}

func (e *JSONResponseSender) createError(err error) ResponseResult {
	return ResponseResult{Message: err.Error(), Data: nil}
}

func (json *JSONResponseSender) SendError(ctx *gin.Context, httpStatus int, err error) {
	errorMsg := json.createError(err)
	ctx.JSON(httpStatus, errorMsg)
}

func (json *JSONResponseSender) createOKRespone(msg string, data interface{}) ResponseResult {
	return ResponseResult{Message: msg, Data: data}
}

func (json *JSONResponseSender) SendOKRespnse(ctx *gin.Context, msg string, data interface{}) {
	response := json.createOKRespone(msg, data)

	ctx.JSON(http.StatusOK, response)
}

func NewJSONResponseSender() ResponseSender {
	return &JSONResponseSender{}
}

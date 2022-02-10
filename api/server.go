package api

import (
	db "bank/db/sqlc"
	"bank/util"

	"github.com/gin-gonic/gin"
)

type Server interface {
	getAccount(ctx *gin.Context)
	createAccount(ctx *gin.Context)
	deleteAccount(ctx *gin.Context)
	Start(address string) error
}

type HttpServer struct {
	store        db.Store
	router       *gin.Engine
	resultSender util.ResponseSender
}

func (server *HttpServer) Start(address string) error {
	return server.router.Run(address)
}

func NewServer(store db.Store, sender util.ResponseSender) *HttpServer {
	router := gin.Default()
	server := &HttpServer{store: store, router: router, resultSender: sender}

	router.POST("accounts", server.createAccount)
	router.PATCH("accounts", server.updateAccountBalance)
	router.DELETE("account/:id", server.deleteAccount)
	router.GET("account/:id", server.getAccount)

	return server
}

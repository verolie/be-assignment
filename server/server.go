package server

import (
	"example/transaction/handler"
	"example/transaction/transaction"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	e := gin.Default()

	registerServer(e)

	e.Run(":8081")
}


func registerServer(e *gin.Engine) {
	// Account Manager Service
	e.POST("/user/login", getUser)
	e.POST("/user/register", createUser)
	e.GET("/user/account/detail/:acct_num", getDetaiUserAccount)
	e.GET("/user/payment/history/:loc_acct", getUserTrnxHist)

	// Transaction
	transactionGroup := e.Group("/transaction")
	transactionGroup.Use(authMiddleware())
	{
		transactionGroup.POST("/send", paymentProcess)
		transactionGroup.POST("/withdraw", withdrawProcess)
	}
	e.GET("/transaction/detail", detailTransaction)
	e.GET("/transaction/detail/:trx_id", detailTransactionParam)
}

func getUser(c *gin.Context) {
	handler.GetUsers(c)
}
func createUser(c *gin.Context) {
	handler.CreateUser(c)
}
func getDetaiUserAccount(c *gin.Context) {
	handler.DetaiUserAccount(c)
}
func getUserTrnxHist(c *gin.Context) {
	handler.UserTrnxHist(c)
}
func paymentProcess(c *gin.Context) {
	transaction.PostingPayment(c)
}

func withdrawProcess(c *gin.Context) {
	transaction.WithdrawProcess(c)
}

func detailTransaction(c *gin.Context) {
	transaction.GetDetailTransaction(c)
}

func detailTransactionParam(c *gin.Context) {
	transaction.GetDetailTransactionParam(c)
}

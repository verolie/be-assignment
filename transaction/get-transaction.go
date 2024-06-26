package transaction

import (
	"context"
	"log"
	"net/http"

	"example/transaction/prisma/db"

	"github.com/gin-gonic/gin"
)


func GetDetailTransaction(c *gin.Context) {
	client := db.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
		return
    }
    defer client.Prisma.Disconnect()

    resp, err := client.TransactionDetail.FindMany().Exec(context.Background())
    if err != nil {
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transactions" , err.Error())))
        return
    }

	c.JSON(http.StatusOK, ResponseDataDetail(resp))
}

func GetDetailTransactionParam(c *gin.Context) {
	client := db.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
		return
    }
    defer client.Prisma.Disconnect()

	trxId := c.Param("trx_id")
    if trxId == "" {
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Transaction ID is required" , "")))
        return
    }

	resp, err := client.TransactionDetail.FindUnique(db.TransactionDetail.TrxID.Equals(trxId)).Exec(context.Background())
    if err != nil {
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transaction details" , "")))
    }

	c.JSON(http.StatusOK, ResponseDataDetail(resp))
}

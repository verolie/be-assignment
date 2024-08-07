package transaction

import (
	"context"
	"log"
	"net/http"

	"example/transaction/prisma/db"

	"github.com/gin-gonic/gin"
)

// GetDetailTransaction retrieves all transaction details from the database
func GetDetailTransaction(c *gin.Context) {
    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err := client.Prisma.Connect(); err != nil {
        // Log fatal error and exit if there's an error connecting to the database
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Retrieve all transaction details from the database
    resp, err := client.TransactionDetail.FindMany().Exec(context.Background())
    if err != nil {
        // Return a 500 Internal Server Error response if there's an error retrieving transactions
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transactions", err.Error())))
        return
    }

    // Return a 200 OK response with the transaction details
    c.JSON(http.StatusOK, ResponseDataDetail(resp))
}

// GetDetailTransactionParam retrieves details for a specific transaction identified by the transaction ID
func GetDetailTransactionParam(c *gin.Context) {
    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err := client.Prisma.Connect(); err != nil {
        // Log fatal error and exit if there's an error connecting to the database
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Retrieve the transaction ID from the URL parameters
    trxId := c.Param("trx_id")
    if trxId == "" {
        // Return a 400 Bad Request response if the transaction ID is not provided
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Transaction ID is required", "")))
        return
    }

    // Retrieve the transaction details for the specified transaction ID from the database
    resp, err := client.TransactionDetail.FindUnique(db.TransactionDetail.TrxID.Equals(trxId)).Exec(context.Background())
    if err != nil {
        // Return a 500 Internal Server Error response if there's an error retrieving the transaction details
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transaction details", "")))
    }

    // Return a 200 OK response with the transaction details
    c.JSON(http.StatusOK, ResponseDataDetail(resp))
}

package handler

import (
	"context"
	"log"
	"net/http"

	"example/transaction/prisma/db"

	"github.com/gin-gonic/gin"
)

func DetaiUserAccount(c *gin.Context) {
    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Get the account number from the URL parameters
    acctNumber := c.Param("acct_num")
    if acctNumber == "" {
        // If the account number is missing, return a 400 Bad Request response
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Transaction ID is required", "")))
        return
    }

    // Retrieve the account details for the given account number
    resp, err := client.AccountDetail.FindMany(db.AccountDetail.AcctNum.Equals(acctNumber)).Exec(context.Background())
    if err != nil {
        // If there's an error retrieving the account details, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transaction details", "")))
        return
    }

    // Return a 200 OK response with the retrieved account details
    c.JSON(http.StatusOK, ResponseDataDetail(resp))
}
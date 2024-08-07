package handler

import (
	"context"
	"log"
	"net/http"

	"example/transaction/prisma/db"

	"github.com/gin-gonic/gin"
)

func UserTrnxHist(c *gin.Context) {
    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Get the location account number from the URL parameters
    locAcct := c.Param("loc_acct")
    if locAcct == "" {
        // If the location account number is missing, return a 400 Bad Request response
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Transaction ID is required", "")))
        return
    }

    // Retrieve the transaction details for the given location account number
    resp, err := client.TransactionDetail.FindMany(db.TransactionDetail.LocAcct.Equals(locAcct)).Exec(context.Background())
    if err != nil {
        // If there's an error retrieving the transaction details, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error retrieving transaction details", err.Error())))
        return
    }

    // Return a 200 OK response with the retrieved transaction details
    c.JSON(http.StatusOK, ResponseDataDetail(resp))
}
package transaction

import (
	"context"
	"example/transaction/model"
	"example/transaction/prisma/db"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WithdrawProcess(c *gin.Context) {
    // Define a variable to store the payment details
    var payment model.DetailTransaction

    // Define a variable to hold any errors
    var err error

    // Set the default status
    status := "Success"

    // Retrieve the email from the context
    email, exists := c.Get("email")
    if !exists {
        // If email is not found in context, return a 400 Bad Request response
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Invalid Email", "")))
        return
    }

    // Bind the JSON request body to the payment variable
    if err = c.BindJSON(&payment); err != nil {
        // If there's an error binding the JSON, return a 400 Bad Request response
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Invalid request body", err.Error())))
        return
    }

    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err = client.Prisma.Connect(); err != nil {
        // If there's an error connecting to the database, log the error and return
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Find the user by email
    user, err := client.User.FindUnique(
        db.User.Email.Equals(email.(string)),
    ).Exec(context.Background())
    if err != nil && err.Error() != "ErrNotFound" || user == nil {
        // If there's an error retrieving the user or the user does not exist, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Email Not Match", err.Error())))
        return
    }

    // Check if the account exists
    if CheckAccount(client, payment.Loc_acct) {
        // Try to update the account if its active type is not "W"
        if actvTyp != "W" {
            UpdateAccount(c, client, payment)
            if status == "Failed" {
                // If updating the account fails, return early
                return
            }
        } else {
            // If the account is already written off, return a 500 Internal Server Error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Account Already Write Off", "")))
            return
        }

        // Insert the transaction details into the database
        _, err = client.TransactionDetail.CreateOne(
            db.TransactionDetail.TrxID.Set(payment.Trx_id),
            db.TransactionDetail.Timestamps.Set(time.Now()),
            db.TransactionDetail.ReceiverPan.Set(""),
            db.TransactionDetail.SenderPan.Set(payment.Sender_pan),
            db.TransactionDetail.ApvCode.Set(payment.Apv_code),
            db.TransactionDetail.TrxTyp.Set(payment.Trx_typ),
            db.TransactionDetail.Amt.Set(payment.Amt),
            db.TransactionDetail.Status.Set(payment.Status),
            db.TransactionDetail.Desc.Set(payment.Desc),
            db.TransactionDetail.AcctDetail.Link(db.AccountDetail.LocAcct.Equals(payment.Loc_acct)),
        ).Exec(context.Background())
        if err != nil {
            // If there's an error inserting the transaction data, return a 500 Internal Server Error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error inserting payment data", err.Error())))
            return
        }
    } else {
        // If the account cannot be found, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Cannot Find Data", "")))
        return
    }

    // Return a 200 OK response indicating the withdrawal was successful
    c.JSON(http.StatusOK, ResponseDataDetail("withdraw success"))
}


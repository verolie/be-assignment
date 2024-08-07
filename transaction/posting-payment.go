package transaction

import (
	"context"
	"example/transaction/model"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"example/transaction/prisma/db"
)

var acctTyp string
var actvTyp, status string
var blncAmt, loanAmt, minLoanPymnt float64

func PostingPayment(c *gin.Context) {
    // Define variables for payment details and error handling
    var payment model.DetailTransaction
    var err error
    status = "Success" // Initialize status

    // Retrieve the email from the context
    email, exists := c.Get("email")
    if !exists {
        // Return a 400 Bad Request response if email is not present in the context
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Email null", err.Error())))
        return
    }

    // Bind the JSON request body to the payment variable
    if err = c.BindJSON(&payment); err != nil {
        // Return a 400 Bad Request response if there is an error with the request body
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Invalid request body", err.Error())))
        return
    }

    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err = client.Prisma.Connect(); err != nil {
        // Log fatal error and return if there's an error connecting to the database
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Retrieve the user from the database using the email from the context
    user, err := client.User.FindUnique(
        db.User.Email.Equals(email.(string)),
    ).Exec(context.Background())
    if (err != nil && err.Error() != "ErrNotFound") || user == nil {
        // Return a 500 Internal Server Error response if there's an error retrieving the user or the user is not found
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Email Not Match", err.Error())))
        return
    }
    
    // Check if the account exists and get its details
    if CheckAccount(client, payment.Loc_acct) {
        // Try to update the account if it exists
        if actvTyp != "W" { // Check if account is not already written off
            UpdateAccount(c, client, payment)
            if status == "Failed" {
                return
            }
        } else {
            // Return a 500 Internal Server Error response if the account is already written off
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Account Already Write Off", err.Error())))
            return
        }

        // Insert the payment transaction details into the database
        _, err = client.TransactionDetail.CreateOne(
            db.TransactionDetail.TrxID.Set(payment.Trx_id),
            db.TransactionDetail.Timestamps.Set(time.Now()),
            db.TransactionDetail.ReceiverPan.Set(payment.Receiver_pan),
            db.TransactionDetail.SenderPan.Set(payment.Sender_pan),
            db.TransactionDetail.ApvCode.Set(payment.Apv_code),
            db.TransactionDetail.TrxTyp.Set(payment.Trx_typ),
            db.TransactionDetail.Amt.Set(float64(payment.Amt)),
            db.TransactionDetail.Status.Set(payment.Status),
            db.TransactionDetail.Desc.Set(payment.Desc),
            db.TransactionDetail.AcctDetail.Link(db.AccountDetail.LocAcct.Equals(payment.Loc_acct)),
        ).Exec(context.Background())
        if err != nil {
            // Return a 500 Internal Server Error response if there's an error inserting the payment data
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error inserting payment data", err.Error())))
            return
        }
    } else {
        // Return a 500 Internal Server Error response if the account cannot be found
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Cannot Find Data", "")))
        return
    }
    
    // Return a 200 OK response indicating the payment was successfully processed
    c.JSON(http.StatusOK, ResponseDataDetail("Success Transaction"))
}

// CheckAccount checks if the account exists and retrieves its details
func CheckAccount(client *db.PrismaClient, Loc_acct string) bool {
    // Retrieve the account details from the database
    accountDetail, err := client.AccountDetail.FindUnique(
        db.AccountDetail.LocAcct.Equals(Loc_acct),
    ).Exec(context.Background())
    if (err != nil && err.Error() != "ErrNotFound") || accountDetail == nil {
        return false // Account not found or error occurred
    }

    // Set global variables to the account details
    acctTyp = accountDetail.AcctTyp
    actvTyp = accountDetail.ActvTyp
    blncAmt = accountDetail.BlncAmt
    loanAmt = accountDetail.LoanAmt
    minLoanPymnt = accountDetail.MinLoanPymnt

    return true // Account exists
}

// UpdateAccount updates the account details based on the transaction type
func UpdateAccount(c *gin.Context, client *db.PrismaClient, payment model.DetailTransaction) {
    // Handle the transaction based on its type
    if payment.Trx_typ == "C" { // Credit transaction
        print("credit") // Debug print statement
        blncAmt -= float64(payment.Amt)
        if blncAmt < 0 {
            if acctTyp == "C" || acctTyp == "PL" {
                // Update loan amount and minimum payment if balance is negative
                loanAmt += float64(payment.Amt)
                minLoanPymnt = (loanAmt + float64(payment.Amt)) * 0.1
            } else {
                // Return a 500 Internal Server Error response if account balance is not sufficient
                c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error Account Balance Not enough", "")))
                status = "Failed"
                return
            }
        }
    } else if payment.Trx_typ == "D" { // Debit transaction
        blncAmt += float64(payment.Amt)
        if acctTyp == "C" || acctTyp == "PL" {
            // Update loan amount and minimum payment
            loanAmt -= float64(payment.Amt)
            minLoanPymnt = (loanAmt + float64(payment.Amt)) * 0.1
            if loanAmt < 0 {
                blncAmt += math.Abs(loanAmt) // Adjust balance if loan amount is negative
            }
        }
    }

    // Update the account details in the database
    _, err := client.AccountDetail.FindMany(db.AccountDetail.LocAcct.Equals(payment.Loc_acct)).Update(
        db.AccountDetail.BlncAmt.Set(blncAmt),
        db.AccountDetail.LoanAmt.Set(loanAmt),
        db.AccountDetail.MinLoanPymnt.Set(minLoanPymnt),
    ).Exec(context.Background())

    if err != nil {
        // Return a 500 Internal Server Error response if there's an error updating the account
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Update"})
        return
    }
}


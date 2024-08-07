package handler

import (
	"context"
	"log"
	"net/http"

	"example/transaction/model"
	"example/transaction/prisma/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
    // Define a variable to store the registration request
    var registerRequest model.RegisUser
    var err error

    // Bind JSON from the request body to the registerRequest variable
    if err := c.BindJSON(&registerRequest); err != nil {
        // If there's an error, return a 400 Bad Request response with the error details
        c.JSON(http.StatusBadRequest, ResponseErrorDetail(CreateErrorResp("Invalid request body", err.Error())))
        return
    }

    // Create a new database client
    client := db.NewClient()
    // Connect to the database
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
        return
    }
    // Ensure the database connection is closed after the function exits
    defer client.Prisma.Disconnect()

    // Check if the account already exists
    status, err := CheckAccount(client, registerRequest.Acct_num)
    if err != nil && err.Error() != "ErrNotFound" {
        // If there's an error that is not "ErrNotFound", return a 500 Internal Server Error response
        println(err)
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error when check account", err.Error())))
        return
    }
    
    // If the account does not exist, proceed to create the user
    if !status {
        // Hash the user's password
        hashPass, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
        if err != nil {
            // If there's an error hashing the password, return a 500 Internal Server Error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("error when hash password", err.Error())))
            return
        }

        // Set the hashed password in the registration request
        registerRequest.Password = string(hashPass)

        // Create the new user in the database
        _, err = client.User.CreateOne(
            db.User.AcctNum.Set(registerRequest.Acct_num),
            db.User.Name.Set(registerRequest.Name),
            db.User.Email.Set(registerRequest.Email),
            db.User.Password.Set(registerRequest.Password),
            db.User.Address.Set(registerRequest.Address),
        ).Exec(context.Background())
        if err != nil && err.Error() != "ErrNotFound" {
            // If there's an error inserting the user data, return a 500 Internal Server Error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error inserting User data", err.Error())))
            return
        }
    }

    // Create the account details in the database
    _, err = client.AccountDetail.CreateOne(
        db.AccountDetail.LocAcct.Set(registerRequest.DetailAccount.Loc_acct),
        db.AccountDetail.PrinPan.Set(registerRequest.DetailAccount.Prin_pan),
        db.AccountDetail.AcctTyp.Set(registerRequest.DetailAccount.Acct_typ),
        db.AccountDetail.ActvTyp.Set(registerRequest.DetailAccount.Actv_typ),
        db.AccountDetail.BlncAmt.Set(registerRequest.DetailAccount.Blnc_amt),
        db.AccountDetail.LoanAmt.Set(registerRequest.DetailAccount.Loan_amt),
        db.AccountDetail.CyccDay.Set(registerRequest.DetailAccount.Cycc_day),
        db.AccountDetail.MinLoanPymnt.Set(registerRequest.DetailAccount.Min_loan_pymnt),
        db.AccountDetail.Acct.Link(
            db.User.AcctNum.Equals(registerRequest.Acct_num),
        ),
    ).Exec(context.Background())

    if err != nil && err.Error() != "ErrNotFound" {
        // If there's an error inserting the payment data, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Error inserting payment data", err.Error())))
        return
    }

    // Return a 200 OK response indicating the user was successfully created
    c.JSON(http.StatusOK, ResponseDataDetail("user succesfully created"))
}

func CheckAccount(client *db.PrismaClient, Acctnum string) (bool, error) {
    // Find a user by account number
    userAcct, err := client.User.FindUnique(
        db.User.AcctNum.Equals(Acctnum),
    ).Exec(context.Background())
    if err != nil && err.Error() != "ErrNotFound" {
        // If there's an error that is not "ErrNotFound", return false and the error
        return false, err
    }
    if userAcct == nil {
        // If the user account does not exist, return false
        return false, nil
    }
    // If the user account exists, return true
    return true, nil
}


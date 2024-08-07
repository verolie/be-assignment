package handler

import (
	"context"
	"example/transaction/model"
	"example/transaction/prisma/db"
	"example/transaction/utils/token"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
    // Define a variable to store the login request
    var loginRequest model.LoginRequest
    // Bind JSON from the request body to the loginRequest variable
    if err := c.BindJSON(&loginRequest); err != nil {
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

    // Find the first user with the matching email
    users, err := client.User.FindFirst(
        db.User.Email.Equals(loginRequest.Email),
    ).Exec(context.Background())

    if err != nil && err.Error() != "ErrNotFound" {
        // If there's an error that is not "ErrNotFound", log the error and return a 500 Internal Server Error response
        log.Fatalf("Error searching for users: %v", err)
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Internal Server Error", err.Error())))
        return
    }

    // If no user is found, return a 404 Not Found response
    if users == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Verify the password
    err = VerifyPassword(loginRequest.Password, users.Password)
    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        // If the password does not match, return a 500 Internal Server Error response with mismatch error
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Mismatch Hash Password", err.Error())))
        return
    }
  
    // Generate a token for the user
    token, err := token.GenerateToken(users.Email)
    if err != nil {
        // If there's an error generating the token, return a 500 Internal Server Error response
        c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Failed Generate Token", err.Error())))
        return 
    }

    // Return a 200 OK response with the generated token
    c.JSON(http.StatusOK, ResponseDataDetail(token))
}

// VerifyPassword compares the provided password with the hashed password
func VerifyPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

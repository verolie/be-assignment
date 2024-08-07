package server

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func authMiddleware() gin.HandlerFunc {
    // Return a middleware function for handling authentication
    return func(c *gin.Context) {
        // Retrieve the JWT token from the Authorization header
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            // If no token is provided, return an error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Authentication Empty", "")))
            return 
        }

        // Parse and validate the JWT token using the secret key
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Return the secret key to validate the token
            return []byte(os.Getenv("API_SECRET")), nil
        })

        // Check for errors in token parsing or validation
        if err != nil || !token.Valid {
            // If there's an error or the token is invalid, return an error response
            c.JSON(http.StatusInternalServerError, ResponseErrorDetail(CreateErrorResp("Token Invalid", "")))
            return
        }
        
        // If the token is valid, extract claims and set user email in the context
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if emailClaim, ok := claims["email"]; ok {
                if email, ok := emailClaim.(string); ok {
                    // Store the email from claims in the context for later use
                    c.Set("email", email)
                    
                    // Proceed to the next handler
                    c.Next()
                    return
                }
            }
        }
    }
}

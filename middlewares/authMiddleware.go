package middlewares

import (
	"lorry-management/utils"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Split the header to get the token part
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}

// package middlewares

// import (
// 	"database/sql"
// 	"lorry-management/db"
// 	"lorry-management/utils"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
// 			c.Abort()
// 			return
// 		}

// 		// Split the header to get the token part
// 		bearerToken := strings.Split(authHeader, " ")
// 		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := bearerToken[1]
// 		claims, err := utils.ValidateToken(tokenString)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
// 			c.Abort()
// 			return
// 		}

// 		email := claims.Email

// 		// Connect to the database
// 		db, err := db.Connect()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database: " + err.Error()})
// 			c.Abort()
// 			return
// 		}
// 		defer db.Close()

// 		var userID int
// 		err = db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&userID)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			} else {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user ID: " + err.Error()})
// 			}
// 			c.Abort()
// 			return
// 		}

// 		c.Set("user_id", userID)
// 		c.Next()
// 	}
// }

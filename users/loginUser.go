package users

import (
	"fmt"
	"lorry-management/db"
	"lorry-management/models"
	"lorry-management/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginUser function handles user login
func LoginUser(c *gin.Context) {
	var loginData models.LoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(loginData.Email)

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var user models.User
	query := `SELECT user_id, name, email, password, role FROM Users WHERE email = $1`
	row := db.QueryRow(query, loginData.Email)
	err = row.Scan(&user.User_id, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}
	var Userresponse struct {
		Name   string `json:"name"`
		UserId int    `json:"id"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	}
	Userresponse.Name = user.Name
	Userresponse.UserId = user.User_id
	Userresponse.Email = user.Email
	Userresponse.Role = user.Role

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": Userresponse, "token": token /*, "token": token*/})
}

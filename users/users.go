package users

import (
	"lorry-management/db"
	"lorry-management/models"
	"lorry-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser function inserts a new user into the database
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password: " + err.Error()})

		return
	}

	// SQL query to insert a new user
	user.Role = "user"
	query := `INSERT INTO Users (name, email, password, role) VALUES ($1, $2, $3, $4)`

	// Execute the insertion query
	_, err = db.Exec(query, user.Name, user.Email, hashedPassword, user.Role)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user and retrieve details: " + err.Error()})

		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User Created successfully"})

}

func GetAllUsers(c *gin.Context) {
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to the database: " + err.Error()})
		return
	}
	defer db.Close()

	// Execute query to get all users
	rows, err := db.Query("SELECT user_id, name, email, role FROM Users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users: " + err.Error()})
		return
	}
	defer rows.Close()

	// Prepare a slice of User structs to hold the data
	var users []models.User

	// Loop through the rows and scan each record into a User object
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.User_id, &user.Name, &user.Email, &user.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user row: " + err.Error()})
			return
		}
		users = append(users, user)
	}

	// Return the slice as JSON response
	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to the database: " + err.Error()})
		return
	}
	defer db.Close()

	// Execute delete statement
	_, err = db.Exec("DELETE FROM Users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User

	// Bind JSON data to User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to the database: " + err.Error()})
		return
	}
	defer db.Close()

	// Optional password hashing (only if the password field is updated)
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password: " + err.Error()})
			return
		}
		user.Password = hashedPassword
	}

	// Update query with the relevant fields
	query := `UPDATE Users SET name = $1, email = $2, role = $3 WHERE id = $4`
	args := []interface{}{user.Name, user.Email, user.Role, userID}

	// Add password to query only if it's provided
	if user.Password != "" {
		query = `UPDATE Users SET name = $1, email = $2, password = $3, role = $4 WHERE id = $5`
		args = []interface{}{user.Name, user.Email, user.Password, user.Role, userID}
	}

	// Execute the update query
	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

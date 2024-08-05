package controllers

import (
	"database/sql"
	"lorry-management/db"
	"lorry-management/models"
	"lorry-management/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetIncomes retrieves all incomes from the database
func GetIncomes(c *gin.Context) {
	var incomes []models.Income

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "SELECT * FROM incomes"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var income models.Income
		if err := rows.Scan(&income.IncomeID, &income.VechileID, &income.Amount, &income.PaymentDate, &income.Status, &income.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		incomes = append(incomes, income)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incomes)
}

// func GetIncomes(c *gin.Context) {
// 	// Retrieve user ID from context
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	var incomes []models.Income

// 	db, err := db.Connect()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer db.Close()

// 	// Query to join incomes and vehicles, filtering by owner_id
// 	query := `
// 		SELECT i.income_id, i.vehicle_id, i.amount, i.payment_date, i.status, i.description
// 		FROM incomes i
// 		INNER JOIN vehicles v ON i.vehicle_id = v.vehicle_id
// 		WHERE v.owner_id = $1
// 	`
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var income models.Income
// 		if err := rows.Scan(&income.IncomeID, &income.VechileID, &income.Amount, &income.PaymentDate, &income.Status, &income.Description); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		incomes = append(incomes, income)
// 	}

// 	if err := rows.Err(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, incomes)
// }

// GetIncomeByID retrieves a single income by its ID
func GetIncomeByID(c *gin.Context) {
	id := c.Param("id")
	var income models.Income

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "SELECT income_id, lorry_id, amount, payment_date, status, COALESCE(description, 'Default Description') AS description FROM incomes WHERE income_id = $1"

	err = db.QueryRow(query, id).Scan(&income.IncomeID, &income.VechileID, &income.Amount, &income.PaymentDate, &income.Status, &income.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, income)
}

// CreateIncome creates a new income record in the database
func CreateIncome(c *gin.Context) {
	var income models.Income
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `INSERT INTO incomes (vehicle_id, amount, payment_date, status, description)
              VALUES ($1, $2, $3, $4, $5) RETURNING income_id`
	err = db.QueryRow(query, income.VechileID, income.Amount, income.PaymentDate, income.Status, &income.Description).Scan(&income.IncomeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"income": income})
}

// UpdateIncome updates an existing income record in the database
func UpdateIncome(c *gin.Context) {
	id := c.Param("id")
	var income models.Income
	if err := c.ShouldBindJSON(&income); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := `
        UPDATE incomes
        SET vehicle_id = $1, amount = $2, payment_date = $3, status = $4, description= $5
        WHERE income_id = $6
        RETURNING income_id, vehicle_id, amount, payment_date, status, description
    `
	err = db.QueryRow(query, income.VechileID, income.Amount, income.PaymentDate, income.Status, income.Description, id).Scan(
		&income.IncomeID, &income.VechileID, &income.Amount, &income.PaymentDate, &income.Status, &income.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income updated successfully", "income": income})
}

// DeleteIncome deletes an income record from the database
func DeleteIncome(c *gin.Context) {
	id := c.Param("id")
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "DELETE FROM incomes WHERE income_id = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Income not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income deleted successfully"})
}

// GetTotalIncome retrieves the total income for a lorry within a specified date range
func GetTotalIncome(c *gin.Context) {
	vehicleID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var query string
	var args []interface{}
	var totalIncome float64

	// Default values for dates

	// If both start and end dates are provided, use them
	if startDate != "" && endDate != "" {

		formattedStartDate, err := utils.Datefromating(startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		formattedEndDate, err := utils.Datefromating(endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query = `SELECT COALESCE(SUM(amount), 0) AS total_income 
		         FROM incomes 
		         WHERE vehicle_id = $1 AND payment_date BETWEEN $2 AND $3`
		args = []interface{}{vehicleID, formattedStartDate, formattedEndDate}
	} else {
		// Default to the current month if dates are not provided
		now := time.Now()
		firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastDay := firstDay.AddDate(0, 1, -1)

		query = `SELECT COALESCE(SUM(amount), 0) AS total_income 
		         FROM incomes 
		         WHERE vehicle_id = $1 AND payment_date BETWEEN $2 AND $3`
		args = []interface{}{vehicleID, firstDay, lastDay}
	}

	err = db.QueryRow(query, args...).Scan(&totalIncome)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	print(totalIncome)
	c.JSON(http.StatusOK, gin.H{"total_income": totalIncome})
}

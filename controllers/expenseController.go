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

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "SELECT * FROM expenses"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.ExpenseID, &expense.VehicleID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.Receipt, &expense.ExpenseDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// func GetExpenses(c *gin.Context) {
// 	// Retrieve user ID from context
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	var expenses []models.Expense

// 	db, err := db.Connect()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer db.Close()

// 	// Query to join expenses and vehicles, filtering by owner_id
// 	query := `
// 		SELECT e.expense_id, e.vehicle_id, e.category_id, e.amount, e.description, e.receipt, e.expense_date
// 		FROM expenses e
// 		INNER JOIN vehicles v ON e.vehicle_id = v.vehicle_id
// 		WHERE v.owner_id = $1
// 	`
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var expense models.Expense
// 		if err := rows.Scan(&expense.ExpenseID, &expense.VehicleID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.Receipt, &expense.ExpenseDate); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		expenses = append(expenses, expense)
// 	}

// 	if err := rows.Err(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, expenses)
// }

func GetExpenseByID(c *gin.Context) {
	id := c.Param("id")
	var expense models.Expense

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "SELECT * FROM expenses WHERE expense_id = $1"

	err = db.QueryRow(query, id).Scan(&expense.ExpenseID, &expense.CategoryID, &expense.Description, &expense.ExpenseDate, &expense.Amount, &expense.VehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `INSERT INTO expenses (vehicle_id, category_id, amount, description, receipt, expense_date)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING expense_id`
	err = db.QueryRow(query, expense.VehicleID, expense.CategoryID, expense.Amount, expense.Description, expense.Receipt, expense.ExpenseDate).Scan(&expense.ExpenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fetch the total expenses for the vehicle
	var totalExpense float64
	totalQuery := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE vehicle_id = $1`
	err = db.QueryRow(totalQuery, expense.VehicleID).Scan(&totalExpense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"expense":       expense,
		"total_expense": totalExpense,
	})
}

func GetTotalExpense(c *gin.Context) {
	vehicleID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	print(vehicleID)
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var query string
	var args []interface{}
	var totalExpense float64

	// Format to desired format
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
		query = `SELECT COALESCE(SUM(amount), 0) AS total_expense FROM expenses WHERE vehicle_id = $1 AND expense_date BETWEEN $2 AND $3`
		args = []interface{}{vehicleID, formattedStartDate, formattedEndDate}
	} else {
		// Calculate the first and last day of the current month
		now := time.Now()
		firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastDay := firstDay.AddDate(0, 1, -1)

		query = `SELECT COALESCE(SUM(amount), 0) AS total_expense FROM expenses WHERE vehicle_id = $1 AND expense_date BETWEEN $2 AND $3`
		args = []interface{}{vehicleID, firstDay, lastDay}
	}

	err = db.QueryRow(query, args...).Scan(&totalExpense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_expense": totalExpense})
}

func UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
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
        UPDATE expenses
        SET vehicle_id = $1, category_id = $2, amount = $3, description = $4, receipt = $5, expense_date = $6
        WHERE expense_id = $7
        RETURNING expense_id, vehicle_id, category_id, amount, description, receipt, expense_date
    `
	err = db.QueryRow(query, expense.VehicleID, expense.CategoryID, expense.Amount, expense.Description, expense.Receipt, expense.ExpenseDate, id).Scan(
		&expense.ExpenseID, &expense.VehicleID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.Receipt, &expense.ExpenseDate)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully", "expense": expense})
}

func DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	query := "DELETE FROM expenses WHERE expense_id = $1"
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

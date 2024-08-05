package controllers

import (
	"database/sql"
	"lorry-management/db"
	"lorry-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateExpenseCategory(c *gin.Context) {
	var category models.ExpenseCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if category.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expense Category Name is required"})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `
        INSERT INTO expense_categories (category_name)
        VALUES ($1)
        RETURNING category_id, category_name
    `
	err = db.QueryRow(query, category.CategoryName).Scan(&category.CategoryID, &category.CategoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense category created successfully", "category": category})
}

func UpdateExpenseCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.ExpenseCategory
	if err := c.ShouldBindJSON(&category); err != nil {
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
        UPDATE expense_categories
        SET category_name = $1
        WHERE category_id = $2
        RETURNING category_id, category_name
    `
	err = db.QueryRow(query, category.CategoryName, id).Scan(&category.CategoryID, &category.CategoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense category updated successfully", "category": category})
}

func DeleteExpenseCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := "DELETE FROM expense_categories WHERE category_id = $1"
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense category deleted successfully"})
}

func GetAllExpenseCategories(c *gin.Context) {
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := "SELECT category_id, category_name FROM expense_categories"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []models.ExpenseCategory
	for rows.Next() {
		var category models.ExpenseCategory
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func GetExpenseCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var category models.ExpenseCategory
	query := "SELECT category_id, category_name FROM expense_categories WHERE category_id = $1"
	err = db.QueryRow(query, id).Scan(&category.CategoryID, &category.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

package controllers

import (
	"lorry-management/db"
	"lorry-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCurrencySettings retrieves the currency settings for a user
func GetCurrencySettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var settings models.CurrencySettings
	query := "SELECT id, user_id, currency_code FROM currency_settings WHERE user_id = $1"
	err = db.QueryRow(query, userID).Scan(&settings.ID, &settings.UserID, &settings.CurrencyCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Currency settings not found"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateCurrencySettings updates the currency settings for a user
func UpdateCurrencySettings(c *gin.Context) {
	var settings models.CurrencySettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}
	settings.UserID = userID.(int)

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `
        INSERT INTO currency_settings (user_id, currency_code)
        VALUES ($1, $2)
        ON CONFLICT (user_id)
        DO UPDATE SET currency_code = $2
        RETURNING id, user_id, currency_code
    `
	err = db.QueryRow(query, settings.UserID, settings.CurrencyCode).Scan(&settings.ID, &settings.UserID, &settings.CurrencyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update currency settings: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Currency settings updated successfully", "settings": settings})
}

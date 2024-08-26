package controllers

import (
	"database/sql"
	"fmt"
	"lorry-management/db"
	"lorry-management/models"
	"lorry-management/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateVehicle(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast email to string"})
		return
	}
	fmt.Println("Retrieved email:", emailStr) // Debugging line

	vehicle.OwnerID, err = utils.GetOwnerID(db, emailStr)

	query := `
		INSERT INTO vehicles (make, model, year, registration_number, capacity, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING vehicle_id, make, model, year, registration_number, capacity, owner_id
	`

	var createdVehicle models.Vehicle
	err = db.QueryRow(query, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.RegistrationNumber, vehicle.Capacity, vehicle.OwnerID).Scan(
		&createdVehicle.VehicleID,
		&createdVehicle.Make,
		&createdVehicle.Model,
		&createdVehicle.Year,
		&createdVehicle.RegistrationNumber,
		&createdVehicle.Capacity,
		&createdVehicle.OwnerID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert Vehicle and retrieve details: " + err.Error()})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vehicle created successfully", "vehicle": createdVehicle})

}

func UpdateVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the lorry ID from the URL parameter
	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Vehicle ID"})
		return
	}

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	query := `
		UPDATE Vehicles
		SET make = $1, model = $2, year = $3, registration_number = $4, capacity = $5
		WHERE vehicle_id = $6
		RETURNING vehicle_id, make, model, year, registration_number, capacity, owner_id
	`

	var updatedVehicle models.Vehicle
	err = db.QueryRow(query, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.RegistrationNumber, vehicle.Capacity, vehicleID).Scan(
		&updatedVehicle.VehicleID,
		&updatedVehicle.Make,
		&updatedVehicle.Model,
		&updatedVehicle.Year,
		&updatedVehicle.RegistrationNumber,
		&updatedVehicle.Capacity,
		&updatedVehicle.OwnerID,
	)

	if err != nil {
		panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle updated successfully", "vehicle": updatedVehicle})
}

func DeleteVehicle(c *gin.Context) {
	// Parse vehicle ID from URL parameter
	vehicleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	// Start database connection
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Prepare the delete query
	query := "DELETE FROM vehicles WHERE vehicle_id = $1 RETURNING vehicle_id"

	// Execute the delete query
	var deletedVehicleID int
	err = db.QueryRow(query, vehicleID).Scan(&deletedVehicleID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle: " + err.Error()})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted successfully", "vehicle_id": deletedVehicleID})
}

// func GetAllVehicles(c *gin.Context) {
// 	email, exists := c.Get("email")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	db, err := db.Connect()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer db.Close()

// 	// Fetch the owner_id based on the email
// 	var ownerID int
// 	err = db.QueryRow("SELECT user_id FROM users WHERE email = $1", email).Scan(&ownerID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve owner ID: " + err.Error()})
// 		return
// 	}

// 	// Fetch vehicles based on the owner_id
// 	query := `
//         SELECT vehicle_id, make, model, year, registration_number, capacity, owner_id
//         FROM vehicles
//         WHERE owner_id = $1
//     `

// 	rows, err := db.Query(query, ownerID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicles: " + err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	var vehicles []models.Vehicle
// 	for rows.Next() {
// 		var vehicle models.Vehicle
// 		err := rows.Scan(&vehicle.VehicleID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.RegistrationNumber, &vehicle.Capacity, &vehicle.OwnerID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan vehicle: " + err.Error()})
// 			return
// 		}
// 		vehicles = append(vehicles, vehicle)
// 	}

// 	if err := rows.Err(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during rows iteration: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
// }

func GetAllVehicles(c *gin.Context) {
	email, exists := c.Get("email")
	print(email)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cast email to string"})
		return
	}
	fmt.Println("Retrieved email:", emailStr) // Debugging line

	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	// Fetch the owner_id based on the email
	var ownerID int
	err = db.QueryRow("SELECT user_id FROM users WHERE email = $1", emailStr).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve owner ID: " + err.Error()})
		return
	}
	fmt.Println("Retrieved ownerID:", ownerID) // Debugging line

	query := `
		SELECT vehicle_id, make, model, year, registration_number, capacity, owner_id
		FROM vehicles
		WHERE owner_id = $1
	`

	rows, err := db.Query(query, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vehicles: " + err.Error()})
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		err := rows.Scan(&vehicle.VehicleID, &vehicle.Make, &vehicle.Model, &vehicle.Year, &vehicle.RegistrationNumber, &vehicle.Capacity, &vehicle.OwnerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan vehicle: " + err.Error()})
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during rows iteration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
}
func FetchExpensesByVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Database connection configuration
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var query string
	var args []interface{}

	// Check if startDate and endDate are provided
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

		query = `
			SELECT expense_id, vehicle_id, category_id, amount, description, receipt, expense_date
			FROM public.expenses
			WHERE vehicle_id = $1 AND expense_date BETWEEN $2 AND $3
			ORDER BY expense_date DESC
		`
		args = []interface{}{vehicleID, formattedStartDate, formattedEndDate}
	} else {
		query = `
			SELECT expense_id, vehicle_id, category_id, amount, description, receipt, expense_date
			FROM public.expenses
			WHERE vehicle_id = $1
			ORDER BY expense_date DESC
		`
		args = []interface{}{vehicleID}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve expenses: %v", err)})
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ExpenseID, &expense.VehicleID, &expense.CategoryID, &expense.Amount, &expense.Description, &expense.Receipt, &expense.ExpenseDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan expense: %v", err)})
			return
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error occurred during rows iteration: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func FetchIncomeByVehicle(c *gin.Context) {
	vehicleID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Database connection configuration
	db, err := db.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var query string
	var args []interface{}

	// Check if startDate and endDate are provided
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

		query = `
			SELECT income_id, vehicle_id, amount, description, payment_date
			FROM public.incomes
			WHERE vehicle_id = $1 AND payment_date BETWEEN $2 AND $3
			ORDER BY payment_date DESC
		`
		args = []interface{}{vehicleID, formattedStartDate, formattedEndDate}
	} else {
		query = `
			SELECT income_id, vehicle_id, amount, description, payment_date
			FROM public.incomes
			WHERE vehicle_id = $1
			ORDER BY payment_date DESC
		`
		args = []interface{}{vehicleID}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve income: %v", err)})
		return
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		err := rows.Scan(&income.IncomeID, &income.VechileID, &income.Amount, &income.Description, &income.PaymentDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan income: %v", err)})
			return
		}
		incomes = append(incomes, income)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error occurred during rows iteration: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incomes": incomes})
}

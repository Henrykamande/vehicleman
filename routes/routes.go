package routes

import (
	"lorry-management/controllers"
	"lorry-management/middlewares"
	"lorry-management/users"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	// Routes for creating owner and property
	r.POST("/user", users.CreateUser)
	r.GET("/user", users.GetAllUsers)
	r.POST("/login", users.LoginUser)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/currency-settings", controllers.GetCurrencySettings)
		protected.PUT("/currency-settings", controllers.UpdateCurrencySettings)
		protected.POST("/create-vehicle", controllers.CreateVehicle)
		protected.PUT("/update-vehicle/:id", controllers.UpdateVehicle)
		protected.DELETE("/delete-vehicle/:id", controllers.DeleteVehicle)
		protected.GET("/vehicle", controllers.GetAllVehicles)

		// expnse Categories
		protected.POST("/create-expensecat", controllers.CreateExpenseCategory)
		protected.PUT("/update-expensecat/:id", controllers.UpdateVehicle)
		protected.DELETE("/delete-expensecat/:id", controllers.DeleteVehicle)
		protected.GET("/vehicle-expense/:id", controllers.FetchExpensesByVehicle)
		protected.GET("/vehicle-income/:id", controllers.FetchIncomeByVehicle)
		protected.GET("/expensecat", controllers.GetAllExpenseCategories)

		//expnse for a specific vehicle route

		protected.POST("/create-expense", controllers.CreateExpense)
		protected.PUT("/update-expense/:id", controllers.UpdateExpense)
		protected.GET("/expense/:id", controllers.GetExpenseByID)
		protected.GET("/totalexpense/:id", controllers.GetTotalExpense)
		protected.GET("/expenses", controllers.GetExpenses)

		// income for specific vehicle
		protected.POST("/create-income", controllers.CreateIncome)
		protected.GET("/incomes", controllers.GetIncomes)
		protected.PUT("/update-income/:id", controllers.UpdateIncome)
		protected.GET("/income/:id", controllers.GetIncomeByID)
		protected.GET("/totalincome/:id", controllers.GetTotalIncome)

	}
	r.Run(":8080")

}

package main

import (
	"fmt"
	"lorry-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	fmt.Println(" +++++++++++++++++++++++++++++++++ step1")
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost"},

		//AllowOrigins:     []string{"https://vehicleman.onrender.com"},

		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	fmt.Println(" +++++++++++++++++++++++++++++++++step2")

	routes.Routers(r)

	//r.POST("/property", routes.CreateProperty)

}

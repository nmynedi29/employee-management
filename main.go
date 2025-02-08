package main

import (
    "employee-management/handlers"
    "employee-management/database"
	"employee-management/models"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    // Initialize the database connection
    database.InitializeDB()

	database.DB.AutoMigrate(&models.Employee{}, &models.Address{})


    // Create Gin router
    r := gin.Default()

    // Employee routes
    r.POST("/employees", handlers.CreateEmployee(database.DB))
    r.GET("/employees", handlers.GetEmployees(database.DB))
    r.GET("/employees/:id", handlers.GetEmployee(database.DB))
    r.PUT("/employees/:id", handlers.UpdateEmployee(database.DB))
    r.DELETE("/employees/:id", handlers.DeleteEmployee(database.DB))


    // Start server
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Server failed to start: ", err)
    }
}

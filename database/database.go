package database

import (
	"employee-management/models" // Import models package
	"gorm.io/driver/mysql"       // Import MySQL driver for GORM
	"gorm.io/gorm"               // Import GORM
	"log"
)

var DB *gorm.DB

// InitializeDB function initializes the database connection
func InitializeDB() {
	// Database connection string (DSN)
	dsn := "root:rootpassword@tcp(mysql:3306)/employee_management?charset=utf8&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.Employee{}, &models.Address{})
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

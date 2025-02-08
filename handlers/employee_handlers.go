package handlers

import (
    "employee-management/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
    "strconv"
)

// CreateEmployee handler to create a new employee 
func CreateEmployee(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var employee models.Employee
        if err := c.ShouldBindJSON(&employee); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Create(&employee).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Employee created successfully", "employee": employee})
    }
}

// Function to convert string to int 
func stringToInt(s string) int {
    num, _ := strconv.Atoi(s)
    return num
}

// GetEmployees retrieves all employees
func GetEmployees(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var employees []models.Employee
        if err := db.Preload("Address").Find(&employees).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
            return
        }

        c.JSON(http.StatusOK, employees)
    }
}

// GetEmployee retrieves a single employee by ID
func GetEmployee(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        var employee models.Employee

        if err := db.Preload("Address").First(&employee, id).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employee"})
            }
            return
        }

        c.JSON(http.StatusOK, employee)
    }
}

func UpdateEmployee(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var employee models.Employee
        id := c.Param("id")
        if err := db.Preload("Address").First(&employee, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
            return
        }
        if err := c.ShouldBindJSON(&employee); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := db.Omit("CreatedAt").Save(&employee).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
            return
        }
        if employee.Address.Street != "" && employee.Address.City != "" {
            employee.Address.EmployeeID = employee.ID
            if err := db.Save(&employee.Address).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
                return
            }
        }
        c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully", "employee": employee})
    }
}





func DeleteEmployee(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        if err := db.Delete(&models.Employee{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
    }
}

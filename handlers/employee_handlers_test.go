package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"employee-management/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//SetupTestDB creates an in-memory SQLite database.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	if err := db.AutoMigrate(&models.Employee{}, &models.Address{}); err != nil {
		t.Fatalf("Failed to migrate models: %v", err)
	}
	return db
}

//SetupRouter initializes a Gin engine.
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/employees", CreateEmployee(db))
	router.GET("/employees", GetEmployees(db))
	router.GET("/employees/:id", GetEmployee(db))
	router.PUT("/employees/:id", UpdateEmployee(db))
	router.DELETE("/employees/:id", DeleteEmployee(db))
	return router
}

// Tests creating an employee with a valid JSON payload.
func TestCreateEmployee_Success(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	employeePayload := map[string]interface{}{
		"name":     "John Doe",
		"position": "Engineer",
		"salary":   75000.0,
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"state":  "NY",
			"zip":    "10001",
		},
	}
	payloadBytes, err := json.Marshal(employeePayload)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response) 
	assert.NoError(t, err)
	assert.Equal(t, "Employee created successfully", response["message"])

	employeeResp, ok := response["employee"].(map[string]interface{})
	assert.True(t, ok, "employee key should contain a JSON object")
	assert.Equal(t, "John Doe", employeeResp["name"])
	assert.Equal(t, "Engineer", employeeResp["position"])
	assert.Equal(t, 75000.0, employeeResp["salary"])

	addr, ok := employeeResp["address"].(map[string]interface{})
	assert.True(t, ok, "address should be a JSON object")
	assert.Equal(t, "123 Main St", addr["street"])
	assert.Equal(t, "New York", addr["city"])
	assert.Equal(t, "NY", addr["state"])
	assert.Equal(t, "10001", addr["zip"])
}

//Tests that invalid JSON returns a 400 error.
func TestCreateEmployee_InvalidJSON(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestGetEmployees tests retrieving all employees.
func TestGetEmployees(t *testing.T) {
	db := setupTestDB(t)

	employee1 := models.Employee{
		Name:     "Alice",
		Position: "Developer",
		Salary:   85000.0,
		Address: models.Address{
			Street: "456 Park Ave",
			City:   "Boston",
			State:  "MA",
			Zip:    "02108",
		},
	}
	employee2 := models.Employee{
		Name:     "Bob",
		Position: "Designer",
		Salary:   65000.0,
		Address: models.Address{
			Street: "789 Elm St",
			City:   "Chicago",
			State:  "IL",
			Zip:    "60616",
		},
	}
	assert.NoError(t, db.Create(&employee1).Error)
	assert.NoError(t, db.Create(&employee2).Error)

	router := setupRouter(db)
	req, _ := http.NewRequest("GET", "/employees", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var employees []models.Employee
	err := json.Unmarshal(resp.Body.Bytes(), &employees) // err declared here using :=
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(employees), 2)
}

//Tests retrieving a single employee by ID.
func TestGetEmployee_Found(t *testing.T) {
	db := setupTestDB(t)
	employee := models.Employee{
		Name:     "Charlie",
		Position: "Manager",
		Salary:   95000.0,
		Address: models.Address{
			Street: "1010 Maple Rd",
			City:   "Seattle",
			State:  "WA",
			Zip:    "98101",
		},
	}
	assert.NoError(t, db.Create(&employee).Error)

	router := setupRouter(db)
	url := "/employees/" + strconv.Itoa(int(employee.ID))
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var emp models.Employee
	err := json.Unmarshal(resp.Body.Bytes(), &emp) // declare err with :=
	assert.NoError(t, err)
	assert.Equal(t, "Charlie", emp.Name)
	assert.Equal(t, "Manager", emp.Position)
}

// Tests that a non-existent employee returns a 404.
func TestGetEmployee_NotFound(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)
	req, _ := http.NewRequest("GET", "/employees/9999", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// Tests updating an employee's details and address.
func TestUpdateEmployee_Success(t *testing.T) {
	db := setupTestDB(t)
	employee := models.Employee{
		Name:     "David",
		Position: "Analyst",
		Salary:   70000.0,
		Address: models.Address{
			Street: "2020 Oak St",
			City:   "San Francisco",
			State:  "CA",
			Zip:    "94102",
		},
	}
	assert.NoError(t, db.Create(&employee).Error)

	router := setupRouter(db)
	updatePayload := map[string]interface{}{
		"name":     "David Updated",
		"position": "Senior Analyst",
		"salary":   80000.0,
		"address": map[string]interface{}{
			"street": "3030 Pine St",
			"city":   "San Jose",
			"state":  "CA",
			"zip":    "95112",
		},
	}
	payloadBytes, err := json.Marshal(updatePayload)
	assert.NoError(t, err)

	url := "/employees/" + strconv.Itoa(int(employee.ID))
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response) // Declare err using :=
	assert.NoError(t, err)
	assert.Equal(t, "Employee updated successfully", response["message"])

	var updatedEmployee models.Employee
	assert.NoError(t, db.Preload("Address").First(&updatedEmployee, employee.ID).Error)
	assert.Equal(t, "David Updated", updatedEmployee.Name)
	assert.Equal(t, "Senior Analyst", updatedEmployee.Position)
	assert.Equal(t, 80000.0, updatedEmployee.Salary)
	assert.Equal(t, "3030 Pine St", updatedEmployee.Address.Street)
	assert.Equal(t, "San Jose", updatedEmployee.Address.City)
	assert.Equal(t, "CA", updatedEmployee.Address.State)
	assert.Equal(t, "95112", updatedEmployee.Address.Zip)
}

func TestDeleteEmployee_Success(t *testing.T) {
	db := setupTestDB(t)
	employee := models.Employee{
		Name:     "Eve",
		Position: "HR",
		Salary:   60000.0,
		Address: models.Address{
			Street: "404 Not Found Rd",
			City:   "Los Angeles",
			State:  "CA",
			Zip:    "90001",
		},
	}
	assert.NoError(t, db.Create(&employee).Error)

	router := setupRouter(db)
	url := "/employees/" + strconv.Itoa(int(employee.ID))
	req, _ := http.NewRequest("DELETE", url, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var deletedEmployee models.Employee
	err := db.First(&deletedEmployee, employee.ID).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

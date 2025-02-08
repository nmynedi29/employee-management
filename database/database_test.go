package database

import (
	"testing"
	"employee-management/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestInitializeDB(t *testing.T) {
	// Create a sqlmock DB.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// SELECT VERSION() as GORM initializes the connection.
	mock.ExpectQuery(`(?i)^SELECT\s+VERSION\(\)`).
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))

	// SELECT SCHEMA_NAME query.
	mock.ExpectQuery(`(?i)^SELECT\s+SCHEMA_NAME\s+from\s+Information_schema\.SCHEMATA\s+where\s+SCHEMA_NAME\s+LIKE\s+\?.*`).
		WithArgs("%", "").
		WillReturnRows(sqlmock.NewRows([]string{"SCHEMA_NAME"}).AddRow("employee_management"))

	// CREATE TABLE query for the addresses table.
	mock.ExpectExec(`(?i)^CREATE TABLE.*addresses`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//GORM DB using the sqlmock connection.
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	// Call AutoMigrate 
	if err := gormDB.AutoMigrate(&models.Address{}); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Verify 
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met: %v", err)
	}
}

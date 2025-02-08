package models

import (
    "github.com/jinzhu/gorm"
)

type Address struct {
    gorm.Model
    Street    string `json:"street"`
    City      string `json:"city"`
    State     string `json:"state"`
    Zip       string `json:"zip"`
	EmployeeID uint `json:"employee_id" gorm:"foreignkey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type Employee struct {
    gorm.Model
    Name      string  `json:"name"`
    Position  string  `json:"position"`
    Salary    float64 `json:"salary"`
    Address   Address `json:"address"` 
}

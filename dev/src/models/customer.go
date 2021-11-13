package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model  `json:"-"` // https://stackoverflow.com/questions/44003152/hide-fields-in-golang-gorm
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Forename    string     `json:"forename"`
	Surname     string     `json:"surname"`
	DateOfBirth time.Time  `json:"dob" gorm:"column:dob"`
}

func NewCustomer(forename, surname string, dob time.Time) *Customer {
	return &Customer{
		Forename:    forename,
		Surname:     surname,
		DateOfBirth: dob,
	}
}

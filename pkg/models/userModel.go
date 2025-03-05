package models

import "gorm.io/gorm"


type User struct {
	gorm.Model
	UserId string `gorm:"primaryKey;unique"`
	FirstName string
	LastName string
	FullName string
	Email string `gorm:"unique"`
	ProfileUrl string
	Trades []Trade `gorm:"foreignKey:UserId;references:UserId"` 
}

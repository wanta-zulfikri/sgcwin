package entities

import (
	"gorm.io/gorm"
)

type (
	Gameonline struct {
		gorm.Model
		UserID      uint
		Name        string `gorm:"type:varchar(150);not null"`
		Description string `gorm:"type:varchar(255);not null"`
		Image       string `gorm:"type:varchar(150);not null"`
		Video       string `gorm:"type:varchar(150);not null"`
		Pdf         string `gorm:"type:varchar(150);not null"`
		Web         string `gorm:"type:varchar(150);not null"`
		Phone       string `gorm:"type:varchar(15);not null"`
	}

	Submission struct {
		ID     			uint `gorm:"primaryKey;autoIncrement;not null"`
		UserID 			uint
		Name   			string `gorm:"type:varchar(255);not null"` 
		Date            string `gorm:"type:varchar(255);not null"`
	}

	Payment struct {
		gorm.Model 
		GameID        uint 
		Description   string   `gorm:"type:varchar(255);not null"` 
		Image         string   `gorm:"type:varchar(70);not null"` 
		Type          string   `gorm:"type:varchar(15);not null"`
		Price         int
		Interval      int 
	}
)

package entities 

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model      `json:"-"` 
		Username          string   `gorm:"type:varchar(30);not null" json:"username,omitempty"`
		Password          string   `gorm:"type:varchar(80);not null" json:"password,omitempty"`    
		IsVerified        bool   `gorm:"not null" json:"-"`
		VerificationCode  string `gorm:"not null" json:"-"`
	} 

	ForgotPass struct {
		Token    string 
		Email    string 
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	LoginReq struct {
		Username string    `json:"username" validate:"required"`
		Password string    `json:"password" validate:"required"`
	} 

	RegisterReq struct {
		Username  string `json:"username" validate:"required"`
		Password  string `json:"password" validate:"required"`

	} 

	UpdateReq struct {
		Id      int 
		Password  string `form:"password"`
		Username  string `form:"username"`
	}
)
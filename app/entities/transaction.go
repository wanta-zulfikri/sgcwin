package entities

import (
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type (
	ReqCharge struct {
		PaymentType string
		Invoice     string
		Total       int
		ItemDetails *[]midtrans.ItemDetails
		CustomerDetails  *midtrans.CustomerDetails
	}
	Carts struct {
		UserID   uint              `gorm:"not null"`
		DeleteAt gorm.DeletedAt    `gorm:"index"`
		User     User
	}
	BillingSchedule struct {
		ID        uint             `gorm:"primaryKey;not null;autoIncrement"`
		Name      string   
		Email     string 
		DeletedAt  gorm.DeletedAt  `gorm:"index"` 
		Date      string           `gorm:"type:timestamp;not null"` 
		Total     int 
	}
	Transaction struct {
		Invoice            string        `gorm:"primaryKey; not null;type:varchar(20)" json:"invoice,omitempty"`
		UserID             uint          `gorm:"not null"` 
		Expire             string        `gorm:"not null"` 
		Total              int           `gorm:"not null"` 
		PaymentCode        string        `gorm:"not null"` 
		PaymentMethod      string        `gorm:"not null"` 
		Status             string        `gorm:"not null"` 
		User               User 
		TransactionItems   []TransactionItems
	}
	TransactionItems struct {
		TransactionInvoice string 
		ItemName           string 
		ItemPrice          int 
	}

	ReqCheckout struct {
		UserID        int     `json:"user_id" validate:"required"` 
		Type          string  `json:"type" validate:"required"` 
		PaymentMethod string  `json:"payment_method" validate:"required"`
	} 
	ResTransaction struct {
		Invoice        string     `json:"invoice"` 
		PaymentMethod  string     `json:"payment_method"` 
		Total          int 		  `json:"total"` 
		PaymentCode    string     `json:"payment_code"` 
		ExpireDate     string     `json:"expire_date"`
	}
	ResGetAllTransaction struct {
		Name       string       `json:"name"` 
		Id         int          `json:"id"`
	}
	ResDetailRegisCart struct {
		Name        string         `json:"name"` 
		Price       int            `json:"price"` 
		Type 		string         `json:"type"` 
		Total       int            `json:"total"`
	}
	ResDetailPayment struct {
		Name        string         `json:"name"` 
		Price       int            `json:"price"`
	}
	ResDetailsHerRegisCart struct {
		OneTime          []ResDetailPayment        `json:"one_time"` 
		Interval         []ResDetailPayment        `json:"interval"` 
		Type             string 
		Total            int 
	}
	ResDetailTransaction struct {
		Invoice          string          `json:"invoice"` 
		PaymentMethod    string          `json:"payment_method"` 
		Total            int             `json:"total"` 
		PaymentCode      string          `json:"payment_code"` 
		Expire           string          `json:"expire"`
	}

)
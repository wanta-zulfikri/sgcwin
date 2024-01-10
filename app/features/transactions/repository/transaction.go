package repository

import (
	entity "mydream_project/app/entities"
	"mydream_project/errorr"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	transaction struct {
		log *logrus.Logger
	}

	TransactionRepo interface {
		CreateTranscation(db  *gorm.DB, data entity.Transaction, typee string) error

	}
)

func NewTransactionRepo(log *logrus.Logger) TransactionRepo {
	return &transaction{log: log}
}

func (t *transaction)  CreateTranscation(db *gorm.DB, data entity.Transaction, typee string ) error {
	return db.Transaction(func(db *gorm.DB) error {
		if err := db.Create(&data).Error; err != nil {
			t.log.Errorf("[ERROR]WHEN Creating Transaction Data, Err: %v", err) 
			return errorr.NewInternal("Internal Server Error")
		}
		status := ""
		if typee == "registration" {
			status = "Send Details Costs Registration"
		} else {
			status = "Send Detail Costs Her-Registration"
		}
		if err := db.Model(&entity.Progress{}).Where("user_id=? AND status != 'Finish' AND status != 'Failed Test Result' AND status != 'Failed File Approved'", data.UserID).Update("status", status).Error; err != nil {
			
		}
	})
}
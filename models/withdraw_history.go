package models

import (
	"time"
)

type WithdrawHistory struct {
	Id        int       `gorm:"primaryKey;" json:"id"`
	UserId    int       `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserId" json:"user"`
	Amount    float64   `gorm:"type:float;" json:"amount"`
	BankName  string    `gorm:"type:varchar(255)" json:"bank_name"`
	AccountNo string    `gorm:"type:varchar(255)" json:"account_no"`
	Status    string    `gorm:"type:varchar(50);default:'Pending Approval'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

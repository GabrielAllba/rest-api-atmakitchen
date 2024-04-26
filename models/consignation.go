package models

type Consignation struct {
	Id          int    `gorm:"primaryKey;" json:"id"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Address     string `gorm:"type:varchar(255)" json:"address"`
	PhoneNumber string `gorm:"type:varchar(255)" json:"phone_number"`
	BankAccount string `gorm:"type:varchar(255)" json:"bank_account"`
	BankNumber  string `gorm:"type:varchar(255)" json:"bank_number"`
}

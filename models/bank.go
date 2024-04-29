package models

type Bank struct {
	Id   int    `gorm:"primaryKey;" json:"id"`
	Name string `gorm:"type:varchar(255)" json:"name"`
	Code string `gorm:"type:varchar(255)" json:"code"`
}

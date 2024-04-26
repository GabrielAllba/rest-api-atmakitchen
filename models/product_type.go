package models

type ProductType struct {
	Id   int    `gorm:"primaryKey;" json:"id"`
	Name string `gorm:"type:varchar(255)" json:"name"`
}

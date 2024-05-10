package models

type BahanProduct struct {
	Id          int     `gorm:"primaryKey;" json:"id"`
	Product     Product `gorm:"foreignKey:ProductId" json:"product"`
	Instruction string  `gorm:"type:varchar(255)" json:"instruction"`
	ProductId   int     `gorm:"index" json:"product_id"`
}
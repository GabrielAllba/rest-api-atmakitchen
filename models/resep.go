package models

type Resep struct {
	Id          int          `gorm:"primaryKey;" json:"id"`
	Product     Product      `gorm:"foreignKey:ProductId" json:"product"`
	Instruction string       `gorm:"type:varchar(255)" json:"instruction"`
	ProductId   int          `gorm:"index" json:"product_id"`
	BahanResep  []BahanResep `gorm:"foreignKey:ResepId" json:"bahan_resep"`
}
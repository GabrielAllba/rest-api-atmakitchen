package models

type BahanResep struct {
	Id        int     `gorm:"primaryKey;" json:"id"`
	BahanId   int     `gorm:"index" json:"bahan_id"`
	Bahan     Bahan   `gorm:"foreignKey:BahanId" json:"bahan"`
	Quantity  float64 `gorm:"type:float;" json:"quantity"`
	Unit      string  `gorm:"type:varchar(255)" json:"unit"`
	Product   Product `gorm:"foreignKey:ProductId" json:"product"`
	ProductId int     `gorm:"index" json:"product_id"`
}

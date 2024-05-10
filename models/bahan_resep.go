package models

type BahanResep struct {
	Id       int     `gorm:"primaryKey;" json:"id"`
	ResepId  int     `gorm:"index" json:"resep_id"`
	Resep    Resep   `gorm:"foreignKey:ResepId" json:"resep"`
	BahanId  int     `gorm:"index" json:"bahan_id"`
	Bahan    Bahan   `gorm:"foreignKey:BahanId" json:"bahan"`
	Quantity float64 `gorm:"type:float;" json:"quantity"`
	Unit     string  `gorm:"type:varchar(255)" json:"unit"`
}

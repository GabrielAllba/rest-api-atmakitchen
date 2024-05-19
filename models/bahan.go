package models

type Bahan struct {
	Id     int     `gorm:"primaryKey;" json:"id"`
	Nama   string  `gorm:"type:varchar(255)" json:"nama"`
	Stok   float64 `gorm:"type:float" json:"stok"`
	Satuan string  `gorm:"type:varchar(255);" json:"satuan"`
}

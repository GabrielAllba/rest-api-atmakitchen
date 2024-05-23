package models

import "time"

type PengeluaranLain struct {
	Id               int       `gorm:"primaryKey;" json:"id"`
	Deskripsi        string    `gorm:"type:varchar(255)" json:"deskripsi"`
	Harga            float64   `gorm:"type:float" json:"harga"`
	Metode           string    `gorm:"type:varchar(255);" json:"metode"`
	TanggalPembelian time.Time `gorm:"type:datetime" json:"tanggal_pembelian"`
}

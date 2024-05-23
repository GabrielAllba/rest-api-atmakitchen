package models

type PengeluaranLain struct {
	Id                 int     `gorm:"primaryKey;" json:"id"`
	Deskripsi          string  `gorm:"type:varchar(255)" json:"deskripsi"`
	Harga              float64 `gorm:"type:float" json:"harga"`
	Metode             string  `gorm:"type:varchar(255);" json:"metode"`
	TanggalPengeluaran string  `gorm:"type:varchar(255)" json:"tanggal_pengeluaran"`
}
